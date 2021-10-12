package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type doesAccountExistScenario struct {
	scenarioName        string
	shouldErrorOccur    bool
	account             *models.MerchantAccount
	shouldCreateAccount bool
	expectedError       error
}

// doesAccountExistTestScenarios returns a set of scenarios to test the check account existence operation
func doesAccountExistTestScenarios() []doesAccountExistScenario {
	return []doesAccountExistScenario{
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
			scenarioName:        "account does not exist",
			shouldErrorOccur:    true,
			account:             GenerateRandomizedAccount(),
			shouldCreateAccount: false,
			expectedError:       service_errors.ErrInvalidInputArguments,
		},
		{
			// failure condition: account does not exist - id(not found)
			scenarioName:        "account does not exist",
			shouldErrorOccur:    true,
			account:             GenerateRandomizedAccountWithRandomId(),
			shouldCreateAccount: false,
			expectedError:       service_errors.ErrAccountDoesNotExist,
		},
	}
}

func TestDoesAccountExistsOperation(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := doesAccountExistTestScenarios()
	for _, scenario := range scenarios {
		var merchantAcct = scenario.account

		if scenario.shouldCreateAccount {
			acct, err := Conn.CreateMerchantAccount(ctx, scenario.account)
			if err != nil {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}

			merchantAcct = acct
		}

		accountExists, err := Conn.CheckAccountExistenceStatus(ctx, merchantAcct.Id)
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
