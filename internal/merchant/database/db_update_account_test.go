package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type updateAccountScenario struct {
	scenarioName        string
	account             *models.MerchantAccount
	shouldCreateAccount bool
	shouldErrorOccur    bool
	expectedError       error
}

func updateAccountScenarios() []updateAccountScenario {
	return []updateAccountScenario{
		{
			// success condition - update existing merchant account
			scenarioName:        "update existing account",
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: true,
			shouldErrorOccur:    false,
			expectedError:       nil,
		},
		{
			// failure condition - update non existing merchant account
			scenarioName:        "update non-existing account",
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: false,
			shouldErrorOccur:    true,
			expectedError:       service_errors.ErrInvalidInputArguments,
		},
		{
			// failure condition - account id is invalid
			scenarioName:        "update invalid account object",
			account:             &models.MerchantAccount{},
			shouldCreateAccount: false,
			shouldErrorOccur:    true,
			expectedError:       service_errors.ErrInvalidInputArguments,
		},
		{
			// failure condition - account does not exist
			account:             GenerateRandomizedAccountWithRandomId(),
			shouldCreateAccount: false,
			shouldErrorOccur:    true,
			expectedError:       service_errors.ErrAccountDoesNotExist,
		},
	}
}

func TestUpdateAccountOperation(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := updateAccountScenarios()
	for _, scenario := range scenarios {
		var merchantAcct = scenario.account

		if scenario.shouldCreateAccount {
			acct, err := Conn.CreateMerchantAccount(ctx, scenario.account)
			if err != nil {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}

			merchantAcct = acct
		}

		// update account email
		updatedEmail := fmt.Sprintf("%s@gmail.com", helper.GenerateRandomString(10))

		merchantAcct.BusinessEmail = updatedEmail
		updatedAcct, err := Conn.UpdateMerchantAccount(ctx, merchantAcct.Id, merchantAcct)
		if err != nil {
			if scenario.shouldErrorOccur {
				assert.Equal(t, err, scenario.expectedError)
			} else {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}
		}

		if scenario.shouldErrorOccur && err == nil {
			t.Errorf("expected error to occur but none did")
		}

		if !scenario.shouldErrorOccur {
			assert.NotNil(t, updatedAcct)
			assert.Equal(t, updatedAcct.BusinessEmail, updatedEmail)
			assert.Equal(t, updatedAcct.Id, merchantAcct.Id)
		}
	}
}
