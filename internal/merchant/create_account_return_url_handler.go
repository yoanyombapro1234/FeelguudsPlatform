package merchant

import (
	"context"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
)

type CreateAccountResponse struct {
	MerchantAccount *models.MerchantAccount `json:"merchant_account"`
}

// CreateAccountReturnUrlHandler godoc
// @Summary Starts the second phase of the merchant account creation process
// @Description starts the second the merchant account creation flow for an end user ... this phase is hit after return url
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (m *AccountComponent) CreateAccountReturnUrlHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing
	stripeConnectedAcctId, err := helper.ExtractStripeConnectedAccountIdFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	stripeAccount, err := m.StripeComponent.GetStripeConnectedAccount(stripeConnectedAcctId, nil)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	acct, err := m.Db.FindMerchantAccountByStripeAccountId(ctx, stripeConnectedAcctId)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !stripeAccount.ChargesEnabled {
		// user didn't finish onboarding ... we provide UI prompts allowing them to finish onboarding through stripe later
		// create the user record in our backend but with onboarding status not set to stripe
		acct.AccountOnboardingState = models.MerchantAccountState_PendingOnboardingCompletion
		acct.AccountOnboardingDetails = models.OnboardingStatus_StripeOnboardingStarted
	}

	// user finished onboarding through stripe, thus we can create the user record in our backend
	if stripeAccount.DetailsSubmitted {
		acct.AccountOnboardingState = models.MerchantAccountState_ActiveAndOnboarded
		acct.AccountOnboardingDetails = models.OnboardingStatus_StripeOnboardingCompleted
	}

	// save the merchant record
	if err := m.Db.SaveAccountRecord(ctx, acct); err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: redirect to merchant account admin page (think about how to do this)
	helper.JSONResponse(w, &CreateAccountResponse{MerchantAccount: acct})
}
