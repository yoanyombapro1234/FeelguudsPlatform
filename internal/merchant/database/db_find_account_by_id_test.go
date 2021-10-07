package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type dbFindAccountByIdScenario struct {
	scenarioName        string
	shouldErrorOccur    bool
	account             *models.MerchantAccount
	shouldCreateAccount bool
	expectedError       error
	deactivateAccount   bool
}

func TestDbFindAccountById(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := getDbFindAccountByIdScenarios()
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

		accountExists, err := Conn.FindMerchantAccountById(ctx, merchantAcct.Id)
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
			assert.True(t, accountExists)
		}
	}
}

// getDbFindAccountByIdScenarios returns a set of scenarios to test the account's existence based on provided id
func getDbFindAccountByIdScenarios() []dbFindAccountByIdScenario {
	testAcct := GenerateRandomizedAccount()
	testAcct.Id = 1000

	return []dbFindAccountByIdScenario{
		{
			// success condition: account exists
			scenarioName:        "account exists",
			shouldErrorOccur:    false,
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: true,
			expectedError:       nil,
		},
		{
			// failure condition: account does not exist - id (0)
			scenarioName:        "account does not exist - id (0)",
			shouldErrorOccur:    true,
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: false,
			expectedError:       service_errors.ErrInvalidInputArguments,
		},
		{
			// failure condition: account does not exist - id (non-existent)
			scenarioName:        "account does not exist - id non-existent",
			shouldErrorOccur:    true,
			account:             testAcct,
			shouldCreateAccount: false,
			expectedError:       service_errors.ErrAccountDoesNotExist,
		},
		{
			// failure condition: account does not exist .. account not active
			scenarioName:        "account does not exists ... account not active",
			shouldErrorOccur:    true,
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: true,
			expectedError:       service_errors.ErrAccountExistButInactive,
			deactivateAccount:   true,
		},
	}
}
