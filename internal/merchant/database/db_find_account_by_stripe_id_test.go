package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type findAccountByStripeIdScenario struct {
	scenarioName        string
	shouldErrorOccur    bool
	account             *models.MerchantAccount
	shouldCreateAccount bool
	expectedError       error
	deactivateAccount   bool
}

// findAccountByStripeIdScenarios returns a set of scenarios to test the account's existence based on provided email
func findAccountByStripeIdScenarios() []findAccountByStripeIdScenario {
	return []findAccountByStripeIdScenario{
		{
			// success condition: account exists
			scenarioName:        "find an account by stripe id that already exists",
			shouldErrorOccur:    false,
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: true,
			expectedError:       nil,
		},
		{
			// failure condition: account does not exist
			scenarioName:        "account does not exist",
			shouldErrorOccur:    true,
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: false,
			expectedError:       service_errors.ErrAccountDoesNotExist,
		},
		{
			// failure condition: account does not exist ... account not active
			scenarioName:        "account does not exists ... account not active",
			shouldErrorOccur:    true,
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: true,
			expectedError:       service_errors.ErrAccountExistButInactive,
			deactivateAccount:   true,
		},
	}
}

func TestFindAccountByStripeIdOperation(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := findAccountByStripeIdScenarios()
	for _, scenario := range scenarios {
		var merchantAcct = scenario.account

		if scenario.shouldCreateAccount {
			acct, err := Conn.CreateMerchantAccount(ctx, scenario.account)
			if err != nil {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}

			if scenario.deactivateAccount {
				ok, err := Conn.DeactivateMerchantAccount(ctx, acct.Id)
				if err != nil {
					t.Errorf("obtained error but not expected - %s", err.Error())
				}

				if !ok {
					t.Errorf("failed to deactivate account")
				}
			}

			merchantAcct = acct
		}

		acct, err := Conn.FindMerchantAccountByStripeAccountId(ctx, merchantAcct.StripeConnectedAccountId)
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
			assert.NotNil(t, acct)
		}
	}
}
