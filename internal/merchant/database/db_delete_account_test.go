package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type dbDeleteAccountScenario struct {
	scenarioName                string
	shouldErrorOccur            bool
	shouldAccountBeCreatedFirst bool
	account                     *models.MerchantAccount
	expectedError               error
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := getDeleteAccountTestScenarios()
	for _, scenario := range scenarios {
		account := scenario.account
		if scenario.shouldAccountBeCreatedFirst {
			acct, err := Conn.CreateMerchantAccount(ctx, account)
			if err != nil {
				t.Errorf("failed to create test account as precondition - %s", err.Error())
			}

			assert.NotNil(t, acct)
			account = acct
		}

		accountDeactivated, err := Conn.DeactivateMerchantAccount(ctx, account.Id)
		if err != nil {
			if scenario.shouldErrorOccur {
				assert.Equal(t, err, scenario.expectedError)
			} else {
				t.Errorf("obtained unexpected error - %s", err.Error())
			}
		}

		if !scenario.shouldErrorOccur {
			assert.True(t, accountDeactivated)
		}
	}

}

// getDeleteAccountTestScenarios returns a set of scenarios to test the delete account operation
func getDeleteAccountTestScenarios() []dbDeleteAccountScenario {
	testAccount := GenerateRandomizedAccount()
	nonExistentAccount := GenerateRandomizedAccount()
	nonExistentAccount.Id = 10000

	return []dbDeleteAccountScenario{
		{
			// success condition - deletion of an existing account
			scenarioName:                "delete existing account",
			shouldErrorOccur:            false,
			shouldAccountBeCreatedFirst: true,
			expectedError:               nil,
			account:                     testAccount,
		},
		{
			// failure condition - invalid account parameters
			scenarioName:                "delete invalid account - invalid account object",
			shouldErrorOccur:            true,
			shouldAccountBeCreatedFirst: false,
			expectedError:               service_errors.ErrInvalidInputArguments,
			account:                     &models.MerchantAccount{},
		},
		{
			// failure condition - deletion of non-existing account
			scenarioName:                "delete non-existent account",
			shouldErrorOccur:            true,
			shouldAccountBeCreatedFirst: false,
			expectedError:               service_errors.ErrAccountDoesNotExist,
			account:                     nonExistentAccount,
		},
	}
}
