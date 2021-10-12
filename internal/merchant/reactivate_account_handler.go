package merchant

import (
	"context"
	"errors"
	"net/http"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
)

// ReactivateMerchantAccountHandler godoc
// @Summary activates a merchant account
// @Description coordinates interactions across multiple services to activate a merchant account
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (m *AccountComponent) ReactivateMerchantAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	id, err := helper.ExtractIDFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	acct, err := m.Db.GetMerchantAccountById(ctx, id, false)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if acct.IsActive {
		return
	}

	sagaSteps, err := m.reactivateMerchantAccountDtxSagaSteps(ctx, acct)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.SagaCoordinater.RunSaga(ctx, "reactivate_merchant_account", sagaSteps...); err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}

// reactivateMerchantAccountDtxSagaSteps returns saga encompassing numerous distributed tx used as part of account re-activation process
func (m *AccountComponent) reactivateMerchantAccountDtxSagaSteps(ctx context.Context, acct *models.MerchantAccount) ([]*saga.Step, error) {
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
		Name: "unlock_merchant_account_distributed_tx",
		Func: func(ctx context.Context) error {
			return m.AuthenticationComponent.UnLockAccount(ctx, authnAccountId)
		},
		CompensateFunc: func(ctx context.Context) error {
			return m.AuthenticationComponent.LockAccount(ctx, authnAccountId)
		},
	}

	activateAcctStep := &saga.Step{
		Name: "activate_merchant_account_op",
		Func: func(ctx context.Context) error {
			if _, err := m.Db.ActivateAccount(ctx, accountId); err != nil {
				return err
			}

			return nil
		},
		CompensateFunc: func(ctx context.Context) error {
			if _, err := m.Db.DeactivateMerchantAccount(ctx, accountId); err != nil {
				return err
			}

			return nil
		},
	}

	sagaSteps = append(sagaSteps, dtxLockAcctStep, activateAcctStep)
	return sagaSteps, nil
}
