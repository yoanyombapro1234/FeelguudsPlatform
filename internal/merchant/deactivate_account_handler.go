package merchant

import (
	"context"
	"errors"
	"net/http"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
)

// DeactivateMerchantAccountHandler godoc
// @Summary Deletes a merchant account
// @Description coordinates interactions across multiple services to delete a merchant account
// @Tags HTTP API
// @Produce html
// @Router / [delete]
// @Success 200 {string} string "OK"
func (m *AccountComponent) DeactivateMerchantAccountHandler(w http.ResponseWriter, r *http.Request) {
	/*
		It is important to note that we never truly ever delete a merchant account. The only thing we do is mark the
		account as deactivated. This enables for quick recovery in case a merchant wants to re-enable the account.

		Due to the distributed nature of this application, the merchant account deletion process requires us to
		both lock the account object from the context of the authentication service and from the context of
		the merchant account which we do in 2 steps. These 2 steps are performed in a saga
	*/

	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	id, err := helper.ExtractIDFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	acct, err := m.Db.GetMerchantAccountById(ctx, id, true)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	sagaSteps, err := m.deactivateMerchantAccountDtxSagaSteps(ctx, acct)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.SagaCoordinater.RunSaga(ctx, "deactivate_merchant_account", sagaSteps...); err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}

// deactivateMerchantAccountDtxSagaSteps returns saga encompassing numerous distributed tx used as part of account deactivation process
func (m *AccountComponent) deactivateMerchantAccountDtxSagaSteps(ctx context.Context, acct *models.MerchantAccount) ([]*saga.Step, error) {
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

	dtxLockAcctStep := &saga.Step{
		Name: "lock_merchant_account_distributed_tx",
		Func: func(ctx context.Context) error {
			return m.AuthenticationComponent.LockAccount(ctx, authnAccountId)
		},
		CompensateFunc: func(ctx context.Context) error {
			return m.AuthenticationComponent.UnLockAccount(ctx, authnAccountId)
		},
	}

	deactivateAcctStep := &saga.Step{
		Name: "deactivate_merchant_account_op",
		Func: func(ctx context.Context) error {
			if _, err := m.Db.DeactivateMerchantAccount(ctx, accountId); err != nil {
				return err
			}

			return nil
		},
		CompensateFunc: func(ctx context.Context) error {
			if _, err := m.Db.ActivateAccount(ctx, accountId); err != nil {
				return err
			}

			return nil
		},
	}

	sagaSteps = append(sagaSteps, dtxLockAcctStep, deactivateAcctStep)
	return sagaSteps, nil
}
