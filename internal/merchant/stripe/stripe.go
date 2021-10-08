package stripe

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/accountlink"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type StripeComponent struct{}

func NewStripeComponent(apiKey string) (*StripeComponent, error) {
	if apiKey == helper.EMPTY {
		return nil, errors.New(fmt.Sprintf("%s - stripe api key cannot be empty", service_errors.ErrInvalidInputArguments.Error()))
	}

	stripe.Key = apiKey
	return &StripeComponent{}, nil
}

// CreateNewStripeConnectedAccount creates a new stripe connected account for a given merchant
func (c *StripeComponent) CreateNewStripeConnectedAccount(merchantAcct *models.MerchantAccount) (string, error) {
	if merchantAcct == nil {
		return helper.EMPTY, errors.New(fmt.Sprintf("%s - stripe param's nil", service_errors.ErrInvalidInputArguments.Error()))
	}

	acctParams := &stripe.AccountParams{
		Params:                stripe.Params{},
		Country:               &merchantAcct.Country,
		DefaultCurrency:       &merchantAcct.DefaultCurrency,
		Email:                 &merchantAcct.BusinessEmail,
		Type:                  stripe.String(string(stripe.AccountTypeStandard)),
		RequestedCapabilities: nil,
	}

	acct, err := account.New(acctParams)
	if err != nil {
		return helper.EMPTY, err
	}

	return acct.ID, nil
}

// CreateNewAccountLink creates a new account link object for a given merchant account
func (c *StripeComponent) CreateNewAccountLink(stripeConnectedAccountId, refreshUrl, returnUrl string) (*stripe.AccountLink, error) {
	params := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeConnectedAccountId),
		RefreshURL: stripe.String(refreshUrl),
		ReturnURL:  stripe.String(returnUrl),
		Type:       stripe.String("account_onboarding"),
	}

	acc, err := accountlink.New(params)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
