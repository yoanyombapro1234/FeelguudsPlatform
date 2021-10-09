package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type activateAccountScenario struct {
	scenarioName            string
	shouldErrorOccur        bool
	account                 *models.MerchantAccount
	shouldDeactivateAccount bool
	shouldCreateAccount     bool
	expectedError           error
}

// activateAccountTestScenarios returns a set of scenarios to test the activate account operation
func activateAccountTestScenarios() []activateAccountScenario {
	return []activateAccountScenario{
		{ // failure condition - attempt to deactivate non-existent account
			scenarioName:            "activate non-existent account",
			shouldErrorOccur:        true,
			account:                 GenerateRandomizedAccountWithRandomId(),
			shouldDeactivateAccount: false,
			shouldCreateAccount:     false,
			expectedError:           service_errors.ErrAccountDoesNotExist,
		},
		{
			// failure condition - attempt to activate an account with invalid id
			scenarioName:            "activate account with invalid id",
			shouldErrorOccur:        true,
			account:                 GenerateRandomizedAccount(),
			shouldDeactivateAccount: false,
			shouldCreateAccount:     false,
			expectedError:           service_errors.ErrInvalidInputArguments,
		},
		{
			// success condition - create a new account and deactivate it
			scenarioName:            "deactivate existing account and attempt re-activation",
			shouldErrorOccur:        false,
			account:                 GenerateRandomizedAccount(),
			shouldDeactivateAccount: true,
			shouldCreateAccount:     true,
			expectedError:           nil,
		},
	}
}

func TestActivateAccountOperation(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := activateAccountTestScenarios()

	var (
		acct *models.MerchantAccount
		err  error
		ok   bool
	)

	for _, scenario := range scenarios {
		acct = scenario.account

		if scenario.shouldCreateAccount {
			if acct, err = Conn.CreateMerchantAccount(ctx, scenario.account); err != nil {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}
		}

		if scenario.shouldDeactivateAccount {
			if ok, err = Conn.DeactivateMerchantAccount(ctx, acct.Id); err != nil {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}

			assert.True(t, ok)
		}

		ok, err := Conn.ActivateAccount(ctx, acct.Id)
		if err != nil {
			if scenario.shouldErrorOccur {
				assert.Equal(t, scenario.expectedError, err)
			} else {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}
		} else {
			if scenario.shouldErrorOccur {
				t.Errorf("expected error %s but none occured", scenario.expectedError.Error())
			}
		}

		if !scenario.shouldErrorOccur {
			assert.True(t, ok)
		}
	}
}
