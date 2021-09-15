package database

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	core_logging "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-logging"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

var (
	db       *Db
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var (
	testBusinessAccount = &models.MerchantAccount{
		Id:                       0,
		Owners:                   nil,
		BusinessName:             "",
		BusinessEmail:            "",
		EmployerId:               0,
		EstimateAnnualRevenue:    "",
		Address:                  nil,
		ItemsOrServicesSold:      nil,
		FulfillmentOptions:       nil,
		ShopSettings:             nil,
		SupportedCauses:          nil,
		Bio:                      "",
		Headline:                 "",
		PhoneNumber:              "",
		Tags:                     nil,
		StripeConnectedAccountId: "",
		StripeAccountId:        0,
		AuthnAccountId:           0,
		AccountOnboardingDetails: 0,
		AccountOnboardingState:   0,
		AccountType:              0,
		IsActive:                 false,
	}
)

func TestMain(m *testing.M) {
	const serviceName string = "test"
	ctx := context.Background()

	// initiate logging client
	logger := InitializeLoggingEngine(ctx)

	params := &helper.DatabaseConnectionParams{
		Host:        host,
		User:        user,
		Password:    password,
		DatabaseName: dbname,
		Port:         port,
	}

	// connect to db
	db, _ = New(ctx, ConnectionInitializationParams{
		ConnectionParams:       params,
		Logger:                 logger,
		MaxConnectionAttempts:  4,
		MaxRetriesPerOperation: 4,
		RetryTimeOut:           3 * time.Second,
		RetrySleepInterval:     50 * time.Millisecond,
	})

	_ = m.Run()
	return
}

// InitializeLoggingEngine initializes a logging object
func InitializeLoggingEngine(ctx context.Context) *zap.Logger {
	// create logging object
	logger := core_logging.New("info")
	return logger.Logger
}

// GenerateRandomId generates a random id over a range
func GenerateRandomId(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

//  ExpectNoErrorOccured ensures no errors occured during the operation
func ExpectNoErrorOccured(t *testing.T, err error, result *models.MerchantAccount) {
	assert.Empty(t, err)
	assert.NotNil(t, result)
}

// ExpectValidAccountObtained ensures we have a valid obtained account
func ExpectValidAccountObtained(t *testing.T, err error, obtainedAccount *models.MerchantAccount, result *models.MerchantAccount) {
	assert.Empty(t, err)
	assert.True(t, obtainedAccount != nil)
	assert.Equal(t, obtainedAccount.BusinessEmail, result.BusinessName)
	assert.Equal(t, obtainedAccount.BusinessName, result.BusinessName)
}

// ExpectInvalidArgumentsError ensure the invalid error is present
func ExpectInvalidArgumentsError(t *testing.T, err error, account *models.MerchantAccount) {
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, service_errors.ErrInvalidInputArguments.Error())
	assert.Nil(t, account)
}

// ExpectAccountAlreadyExistError ensures the account already exist error is present
func ExpectAccountAlreadyExistError(t *testing.T, err error, createdAccount *models.MerchantAccount) {
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, service_errors.ErrAccountAlreadyExist.Error())
	assert.Nil(t, createdAccount)
}

// ExpectAccountDoesNotExistError ensures the account does not exist error is present
func ExpectAccountDoesNotExistError(t *testing.T, err error, createdAccount *models.MerchantAccount) {
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, service_errors.ErrAccountDoesNotExist.Error())
	assert.Nil(t, createdAccount)
}

// ExpectCannotUpdatePasswordError ensure the invalid error is present
func ExpectCannotUpdatePasswordError(t *testing.T, err error, account *models.MerchantAccount) {
	assert.NotEmpty(t, err)
	assert.EqualError(t, err, service_errors.ErrCannotUpdatePassword.Error())
	assert.Nil(t, account)
}

// GenerateRandomizedAccount generates a random account
func GenerateRandomizedAccount() *models.MerchantAccount {
	randStr := helper.GenerateRandomString(150)
	account := testBusinessAccount
	account.BusinessName = account.BusinessEmail + randStr
	account.BusinessName = account.BusinessName + randStr
	return account
}
