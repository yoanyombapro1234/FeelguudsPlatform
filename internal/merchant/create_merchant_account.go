package merchant

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/stripe/stripe-go/account"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type CreateMerchantAccountRequest struct {
	MerchantAccount   *models.MerchantAccount `json:"merchant_account"`
	Password          string                  `json:"password"`
	ConfirmedPassword string                  `json:"confirmed_password"`
}

type CreateMerchantAccountResponse struct {
	MerchantAccount *models.MerchantAccount `json:"merchant_account"`
}

// CreateMerchantAccountHandler godoc
// @Summary Starts the first phase of the merchant account creation process
// @Description starts the merchant account creation flow for an end user
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (m *MerchantAccountComponent) CreateMerchantAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing
	var (
		req CreateMerchantAccountRequest
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
	newAcct, err := m.CreateMerchantAccount(ctx, acct)
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

// CreateMerchantAccountRefreshUrlHandler godoc
// @Summary Serves as the handler for the refresh url used throughout the merchant account onboarding process
// @Description resets the onboarding process for a given user
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (m *MerchantAccountComponent) CreateMerchantAccountRefreshUrlHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing
	connectedAcctId, err := helper.ExtractStripeConnectedAccountIdFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshUrl := fmt.Sprintf("%s/%s", m.BaseRefreshUrl, connectedAcctId)
	returnUrl := fmt.Sprintf("%s/%s", m.BaseReturnUrl, connectedAcctId)

	// pull the merchant account from the database
	acct, err := m.Db.FindMerchantAccountByStripeAccountId(ctx, connectedAcctId)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	stripeRedirectUri, err := m.getStripeRedirectURI(acct, connectedAcctId, refreshUrl, returnUrl)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, stripeRedirectUri, http.StatusOK)
}

// CreateMerchantAccountReturnUrlHandler godoc
// @Summary Starts the second phase of the merchant account creation process
// @Description starts the second the merchant account creation flow for an end user ... this phase is hit after return url
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (m *MerchantAccountComponent) CreateMerchantAccountReturnUrlHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing
	stripeConnectedAcctId, err := helper.ExtractStripeConnectedAccountIdFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	stripeAccount, err := account.GetByID(stripeConnectedAcctId, nil)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// pull the merchant account from the database
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

	helper.JSONResponse(w, &CreateMerchantAccountResponse{MerchantAccount: acct})
}

// CreateMerchantAccount creates a merchant account and saves record in database
func (m MerchantAccountComponent) CreateMerchantAccount(ctx context.Context, acct *models.MerchantAccount) (*models.MerchantAccount,
	error) {
	// TODO: start span
	// TODO: emit metrics
	if err := validateMerchantAccount(acct); err != nil {
		return nil, err
	}

	newAcct, err := m.Db.CreateMerchantAccount(ctx, acct)
	if err != nil {
		return nil, err
	}

	return newAcct, nil
}

// getStripeRedirectURI returns the stripe redirect URI enabling end user to start onboarding flow through stripe
func (m MerchantAccountComponent) getStripeRedirectURI(merchantAccount *models.MerchantAccount, connectedAcctId, refreshUrl,
	returnUrl string) (string, error) {

	merchantAccount.StripeConnectedAccountId = connectedAcctId
	link, err := m.StripeComponent.CreateNewAccountLink(merchantAccount.StripeConnectedAccountId, "", "")
	if err != nil {
		return helper.EMPTY, err
	}

	return link.URL, nil
}

// getConnectedAccountId returns the stripe connected account id
func (m MerchantAccountComponent) getConnectedAccountId(merchantAccount *models.MerchantAccount) (string, error) {
	connectedAcctId, err := m.StripeComponent.CreateNewStripeConnectedAccount(merchantAccount)
	if err != nil {
		return helper.EMPTY, err
	}
	return connectedAcctId, nil
}

// validateMerchantAccount validates a merchant account and ensures all required fields are present
func validateMerchantAccount(acc *models.MerchantAccount) error {
	if acc == nil {
		return errors.New(fmt.Sprintf("%s - nil merchant account object", service_errors.ErrInvalidInputArguments.Error()))
	}

	if acc.BusinessEmail == helper.EMPTY {
		return errors.New("merchant account business email cannot be empty")
	}

	if acc.BusinessName == helper.EMPTY {
		return errors.New("merchant account business name cannot be empty")
	}

	if acc.EmployerId == 0 {
		return errors.New("merchant account employer id cannot be empty")
	}

	if acc.AuthnAccountId == 0 {
		return errors.New("merchant account authn id cannot be empty")
	}

	return nil
}
