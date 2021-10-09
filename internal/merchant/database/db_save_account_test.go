package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

type saveAccountScenario struct {
	scenarioName        string
	shouldErrorOccur    bool
	account             *models.MerchantAccount
	expectedError       error
	shouldCreateAccount bool
}

// saveAccountTestScenarios returns a set of scenarios to test the save account operation
func saveAccountTestScenarios() []saveAccountScenario {
	return []saveAccountScenario{
		{
			// success condition - save a new merchant account
			scenarioName:        "save new merchant account",
			shouldErrorOccur:    false,
			account:             GenerateRandomizedAccount(),
			expectedError:       nil,
			shouldCreateAccount: true,
		},
		{
			// failure condition - attempt to save invalid account (nil)
			scenarioName:     "save invalid account object",
			shouldErrorOccur: true,
			account:          nil,
			expectedError:    service_errors.ErrInvalidInputArguments,
		},
		{
			// failure condition - attempt to create invalid account (empty account object)
			scenarioName:     "save invalid account object",
			shouldErrorOccur: true,
			account:          &models.MerchantAccount{},
			expectedError:    service_errors.ErrMisconfiguredIds,
		},
	}
}

func TestSaveAccountOperation(t *testing.T) {
	ctx := context.Background()
	SetupTestDbConn()

	scenarios := saveAccountTestScenarios()
	for _, scenario := range scenarios {
		acct := scenario.account
		if scenario.shouldCreateAccount {
			newAcct, err := Conn.CreateMerchantAccount(ctx, scenario.account)
			if err != nil {
				t.Errorf("failed to create test account as precondition - %s", err.Error())
			}

			assert.NotNil(t, newAcct)

			// update a random field after account is successfully created
			newAcct.BusinessEmail = helper.GenerateRandomString(40)
			acct = newAcct
		}

		if err := Conn.SaveAccountRecord(ctx, acct); err != nil {
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
			// attempt to get the account from the database
			if acct != nil {
				accountFound, err := Conn.FindMerchantAccountByEmail(ctx, acct.BusinessEmail)
				if err != nil {
					t.Errorf("obtained error but not expected - %s", err.Error())
				}

				assert.True(t, accountFound)
			}
		}
	}
}
