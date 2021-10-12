package stripe

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type Interface interface {
	// CreateNewStripeConnectedAccount invokes the stripe API to create a new connected account
	// and returns the connected account ID
	//
	// A stripe connected account enables our merchant to accept payments and move funds to their bank account.
	// Connected accounts represent our user in Stripe’s API and help facilitate the collection of onboarding
	// requirements so Stripe can verify the user’s identity
	CreateNewStripeConnectedAccount(merchantAcct *models.MerchantAccount) (string, error)

	// CreateNewAccountLink invokes the stripe API an returns a set of redirect links which are crucial for
	// merchant the account onboarding process
	CreateNewAccountLink(stripeConnectedAccountId, refreshUrl, returnUrl string) (*stripe.AccountLink, error)
}

type Component struct{}

var _ Interface = (*Component)(nil)

func NewStripeComponent(apiKey string) (*Component, error) {
	if apiKey == helper.EMPTY {
		return nil, errors.New(fmt.Sprintf("%s - stripe api key cannot be empty", service_errors.ErrInvalidInputArguments.Error()))
	}

	stripe.Key = apiKey
	return &Component{}, nil
}

// CreateNewStripeConnectedAccount creates a new stripe connected account for a given merchant
func (c *Component) CreateNewStripeConnectedAccount(merchantAcct *models.MerchantAccount) (string, error) {
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
func (c *Component) CreateNewAccountLink(stripeConnectedAccountId, refreshUrl, returnUrl string) (*stripe.AccountLink, error) {
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

// GetStripeConnectedAccount returns a stripe connected account record
func (c *Component) GetStripeConnectedAccount(acctId string, params *stripe.AccountParams) (*stripe.Account, error){
	if acctId == helper.EMPTY {
		return nil, service_errors.ConcatenateErrorMessages(service_errors.ErrInvalidInputArguments.Error(), "empty stripe connected account ID")
	}

	return account.GetByID(acctId, nil)
}
