package merchant

import (
	"context"
	"errors"
	"net/http"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
)

type UpdateAccountRequest struct {
	MerchantAccount *models.MerchantAccount `json:"merchant_account"`
}

type UpdateMerchantAccountResponse struct {
	MerchantAccount *models.MerchantAccount `json:"merchant_account"`
}

func (m *MerchantAccountComponent) UpdateMerchantAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	/*
				The merchant account update process comprises numerous steps. Upon obtaining a request to update
				a new merchant account, we perform local validations. We then attempt to obtain the account record
		        from the backend database. We do this to cross check if the merchant is attempting to update
		        their email address. if so we invoke the authentication service to update the record.

				Upon acquiring a successful response, we save the updated merchant record in our local backend database.
				It is important to note that the above steps are executed as a set of distributed tx hence we leverage sagas
	*/

	// TODO: emit metrics and add distributed tracing
	var (
		req UpdateAccountRequest
	)

	err := helper.DecodeJSONBody(w, r, &req)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := helper.ExtractIDFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedAcct := req.MerchantAccount
	if err := m.validateMerchantAccount(updatedAcct); err == nil {
		helper.ErrorResponse(w, "invalid merchant account object passed", http.StatusBadRequest)
		return
	}

	// obtain the old account record from db
	oldAcct, err := m.Db.GetMerchantAccountById(ctx, id, true)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// we check is email is being updated and perform the update operation as a distributed tx
	if oldAcct.BusinessEmail != updatedAcct.BusinessEmail {
		oldEmail := oldAcct.BusinessEmail
		newEmail := updatedAcct.BusinessEmail

		// execute distributed tx
		sagaSteps, err := m.updateMerchantAccountDtxSagaSteps(ctx, updatedAcct, newEmail, oldEmail)
		if err != nil {
			helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := m.SagaCoordinater.RunSaga(ctx, "update_merchant_account", sagaSteps...); err != nil {
			helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		helper.JSONResponse(w, UpdateMerchantAccountResponse{updatedAcct})
		return
	}

	// save the record locally
	account, err := m.Db.UpdateMerchantAccount(ctx, id, updatedAcct)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.JSONResponse(w, UpdateMerchantAccountResponse{account})
}

// updateMerchantAccountDtxSagaSteps returns saga encompassing numerous distributed tx used as part of account update process
func (m *MerchantAccountComponent) updateMerchantAccountDtxSagaSteps(ctx context.Context, acct *models.MerchantAccount, newEmail,
	oldEmail string) ([]*saga.Step,
	error) {
	var (
		sagaSteps      = make([]*saga.Step, 0)
		authnAccountId uint32
		accountId      uint64
	)

	if acct == nil {
		return nil, errors.New("invalid input arguments - merchant account object cannot be nil")
	}

	if acct.AuthnAccountId == 0 {
		return nil, errors.New("invalid input arguments - authn account id cannot be 0")
	}

	if acct.Id == 0 {
		return nil, errors.New("invalid input arguments - account id cannot be 0")
	}

	authnAccountId = uint32(acct.AuthnAccountId)
	accountId = acct.Id

	dtxUpdateAcctStep := &saga.Step{
		Name: "update_merchant_account_distributed_tx",
		Func: func(ctx context.Context) error {
			return m.AuthenticationComponent.UpdateAccount(ctx, authnAccountId, newEmail)
		},
		CompensateFunc: func(ctx context.Context) error {
			return m.AuthenticationComponent.UpdateAccount(ctx, authnAccountId, oldEmail)
		},
	}

	updateAcctStep := &saga.Step{
		Name: "update_merchant_account_op",
		Func: func(ctx context.Context) error {
			acct.BusinessEmail = newEmail
			if _, err := m.Db.UpdateMerchantAccount(ctx, accountId, acct); err != nil {
				return err
			}
			return nil
		},
		CompensateFunc: func(ctx context.Context) error {
			acct.BusinessEmail = oldEmail
			if _, err := m.Db.UpdateMerchantAccount(ctx, accountId, acct); err != nil {
				return err
			}

			return nil
		},
	}

	sagaSteps = append(sagaSteps, dtxUpdateAcctStep, updateAcctStep)
	return sagaSteps, nil
}
