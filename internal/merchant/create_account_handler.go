package merchant

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
)

type CreateAccountRequest struct {
	MerchantAccount   *models.MerchantAccount `json:"merchant_account"`
	Password          string                  `json:"password"`
	ConfirmedPassword string                  `json:"confirmed_password"`
}

// CreateAccountHandler godoc
// @Summary Starts the first phase of the merchant account creation process
// @Description starts the merchant account creation flow for an end user
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (m *MerchantAccountComponent) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	/*
			The merchant account creation process comprises numerous steps. Upon obtaining a request to create
			a new merchant account, we perform local validations. Then we invoke the authentication service
			and pass the account's credentials which are comprised of a password and an email.

			upon successful record creation from the context of the authentication service and acquisition of the
			record's id, we update the merchant account record, set the onboarding state and store it internally.

			Once the account record has been successfully saved in the database which the merchant component owns,
		    we commence the stripe onboarding flow. We create an account record on stripe's end and obtain a connected account
			id of which we utilize to generate a redirect URI. The redirect URI is used to enable to client to continue
			the onboarding flow via stripe.
	*/

	// TODO: emit metrics and add distributed tracing
	var (
		req CreateAccountRequest
	)

	err := helper.DecodeJSONBody(w, r, &req)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	acct := req.MerchantAccount
	password := req.Password
	confirmedPassword := req.ConfirmedPassword
	if acct == nil {
		helper.ErrorResponse(w, "invalid merchant account object passed", http.StatusBadRequest)
		return
	}

	if password != confirmedPassword {
		helper.ErrorResponse(w, "password and confirmed password must match", http.StatusBadRequest)
		return
	}

	// call the authentication service and create an account record from its context
	email := acct.BusinessEmail
	authnId, err := m.AuthenticationComponent.CreateAccount(ctx, email, password, false)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// we first create the merchant account record internally and set the onboarding status to not yet started
	// then we generate a stripe redirect URI and commence phase 2 of onboarding
	acct.AuthnAccountId = uint64(authnId)
	acct.AccountOnboardingState = models.MerchantAccountState_PendingOnboardingCompletion
	acct.AccountOnboardingDetails = models.OnboardingStatus_FeelGuudOnboardingStarted
	newAcct, err := m.Db.CreateMerchantAccount(ctx, acct)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	connectedAcctId, err := m.getConnectedAccountId(newAcct)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshUrl := fmt.Sprintf("%s/%s", m.BaseRefreshUrl, connectedAcctId)
	returnUrl := fmt.Sprintf("%s/%s", m.BaseReturnUrl, connectedAcctId)

	stripeRedirectUri, err := m.getStripeRedirectURI(acct, connectedAcctId, refreshUrl, returnUrl)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, stripeRedirectUri, http.StatusOK)
}
