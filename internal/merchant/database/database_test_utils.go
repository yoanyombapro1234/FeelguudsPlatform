package database

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
)

const (
	DefaultMaxConnectionAttempts  int           = 3
	DefaultMaxRetriesPerOperation int           = 3
	DefaultRetryTimeout           time.Duration = 50 * time.Millisecond
	DefaultRetrySleepInterval     time.Duration = 25 * time.Millisecond
)

var (
	Conn                *Db
	Port                int    = 6000
	Host                string = "localhost"
	User                string = "merchant_component"
	Password            string = "merchant_component"
	Dbname              string = "merchant_component"
	DefaultDbConnParams        = helper.DatabaseConnectionParams{
		Host:         Host,
		User:         User,
		Password:     Password,
		DatabaseName: Dbname,
		Port:         Port,
	}

	DefaultConnInitializationParams = ConnectionInitializationParams{
		ConnectionParams:       &DefaultDbConnParams,
		Logger:                 zap.L(),
		MaxConnectionAttempts:  DefaultMaxConnectionAttempts,
		MaxRetriesPerOperation: DefaultMaxRetriesPerOperation,
		RetryTimeOut:           DefaultRetryTimeout,
		RetrySleepInterval:     DefaultRetrySleepInterval,
	}
)

// SetupTestDbConn sets up a database connection to the test db node
func SetupTestDbConn() {
	ctx := context.Background()
	// setup database connection before tests
	Conn, _ = New(ctx, &DefaultConnInitializationParams)
}

// GenerateRandomId generates a random id over a range
func GenerateRandomId(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func GenerateRandomizedAccountWithRandomId() *models.MerchantAccount {
	acct := GenerateRandomizedAccount()
	acct.Id = uint64(GenerateRandomId(10000, 3000000))
	return acct
}

// GenerateRandomizedAccount generates a random account
func GenerateRandomizedAccount() *models.MerchantAccount {
	randStr := helper.GenerateRandomString(5)

	return &models.MerchantAccount{
		Owners: []*models.Owner{
			{
				FirstName: helper.GenerateRandomString(8),
				LastName:  helper.GenerateRandomString(8),
				Email:     fmt.Sprintf("%s@gmail.com", randStr),
				Country:   "USA",
			},
		},
		BusinessName:          randStr,
		BusinessEmail:         fmt.Sprintf("%s@gmail.com", randStr),
		EmployerId:            13450,
		EstimateAnnualRevenue: "1000000",
		Address: &models.Address{
			Address:   "340 Clifton Pl",
			Unit:      "3B",
			ZipCode:   "10013",
			City:      "Brooklyn",
			State:     "NYC",
			Longitude: "40.7131° N",
			Lattitude: "74.0338° W",
		},
		ItemsOrServicesSold: []*models.ItemSold{
			{
				Type: models.ItemSold_SERVICES,
			},
			{
				Type: models.ItemSold_PHYSICAL_ITEMS,
			},
		},
		FulfillmentOptions: []models.FulfillmentOptions{
			models.FulfillmentOptions_SHIP_ITEMS,
			models.FulfillmentOptions_ALLOW_DELIVERY,
		},
		ShopSettings: &models.Settings{
			PaymentDetails: &models.Settings_PaymentDetails{
				AcceptableCreditCardTypes: []models.Settings_PaymentDetails_CreditCardBrand{
					models.Settings_PaymentDetails_VISA,
					models.Settings_PaymentDetails_DISCOVER,
				},
				PrimaryCurrencyCode: models.Settings_PaymentDetails_USD,
				EnabledCurrencyCodes: []models.Settings_PaymentDetails_CurrencyCode{
					models.Settings_PaymentDetails_USD,
					models.Settings_PaymentDetails_GBP,
				},
				SupportedDigitalWallets: []models.Settings_PaymentDetails_DigitalWallets{
					models.Settings_PaymentDetails_APPLE_PAY,
					models.Settings_PaymentDetails_GOOGLE_PAY,
				},
			},
			ShopPolicy:     nil,
			PrivacyPolicy:  nil,
			ReturnPolicy:   nil,
			ShippingPolicy: nil,
		},
		SupportedCauses: []models.Causes{
			models.Causes_EDUCATION,
		},
		Bio:                      "",
		Headline:                 "Creating a better online shopping experience for you",
		PhoneNumber:              "551-778-1002",
		Tags:                     nil,
		StripeConnectedAccountId: helper.GenerateRandomString(15),
		StripeAccountId:          100,
		AuthnAccountId:           40,
		AccountOnboardingDetails: models.OnboardingStatus_OnboardingNotStarted,
		AccountOnboardingState:   models.MerchantAccountState_PendingOnboardingCompletion,
		AccountType:              models.MerchantAccountType_Company,
		IsActive:                 true,
	}
}
