package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type createAccountScenario struct {
	scenarioName     string
	shouldErrorOccur bool
	account          *models.MerchantAccount
	expectedError    error
}

// createAccountTestScenarios returns a set of scenarios to test the create account operation
func createAccountTestScenarios() []createAccountScenario {
	return []createAccountScenario{
		{
			// success condition - create a new merchant account
			scenarioName:     "create new merchant account",
			shouldErrorOccur: false,
			account:          GenerateRandomizedAccount(),
			expectedError:    nil,
		},
		{
			// failure condition - attempt to create duplicate merchant account
			scenarioName:     "create duplicate merchant account",
			shouldErrorOccur: true,
			account:          GenerateRandomizedAccount(),
			expectedError:    service_errors.ErrAccountAlreadyExist,
		},
		{
			// failure condition - attempt to create invalid account (nil)
			scenarioName:     "create invalid account object",
			shouldErrorOccur: true,
			account:          nil,
			expectedError:    service_errors.ErrInvalidAccount,
		},
		{
			// failure condition - attempt to create invalid account (empty account object)
			scenarioName:     "create invalid account object",
			shouldErrorOccur: true,
			account:          &models.MerchantAccount{},
			expectedError:    service_errors.ErrMisconfiguredIds,
		},
	}
}

func TestCreateAccountOperation(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := createAccountTestScenarios()
	for _, scenario := range scenarios {
		account, err := Conn.CreateMerchantAccount(ctx, scenario.account)
		if err != nil {
			if scenario.shouldErrorOccur {
				assert.Equal(t, err, scenario.expectedError)
			} else {
				t.Errorf("obtained error but not expected - %s", err.Error())
			}
		} else {
			if scenario.shouldErrorOccur {
				t.Errorf("expected error %s but none occured", scenario.expectedError.Error())
			}
		}

		if !scenario.shouldErrorOccur {
			assert.NotNil(t, account)
		}
	}
}
