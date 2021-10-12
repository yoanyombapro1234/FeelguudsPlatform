package merchant

import (
	"errors"
	"fmt"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

// getStripeRedirectURI returns the stripe redirect URI enabling end user to start onboarding flow through stripe
func (m AccountComponent) getStripeRedirectURI(merchantAccount *models.MerchantAccount, connectedAcctId, refreshUrl,
	returnUrl string) (string, error) {

	merchantAccount.StripeConnectedAccountId = connectedAcctId
	link, err := m.StripeComponent.CreateNewAccountLink(merchantAccount.StripeConnectedAccountId, refreshUrl, returnUrl)
	if err != nil {
		return helper.EMPTY, err
	}

	return link.URL, nil
}

// getConnectedAccountId returns the stripe connected account id
func (m AccountComponent) getConnectedAccountId(merchantAccount *models.MerchantAccount) (string, error) {
	connectedAcctId, err := m.StripeComponent.CreateNewStripeConnectedAccount(merchantAccount)
	if err != nil {
		return helper.EMPTY, err
	}
	return connectedAcctId, nil
}

// validateMerchantAccount validates a merchant account and ensures all required fields are present
func (m AccountComponent) validateMerchantAccount(acc *models.MerchantAccount) error {
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
