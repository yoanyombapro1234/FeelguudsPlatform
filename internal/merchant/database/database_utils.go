package database

import (
	"context"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
)

const (
	EMPTY = ""
)

// ValidateAccount performs various account level validations
func (db *Db) ValidateAccount(ctx context.Context, account *models.MerchantAccountORM) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err.Error())
		return err
	}

	err = db.ValidateAccountIds(ctx, account)
	if err != nil {
		db.Logger.Error(err.Error())
		return err
	}

	err = db.ValidateAccountParameters(ctx, account)
	if err != nil {
		db.Logger.Error(err.Error())
		return err
	}

	return nil
}

// ValidateAccountParameters validates account params.
func (db *Db) ValidateAccountParameters(ctx context.Context, account *models.MerchantAccountORM) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err.Error())
		return err
	}

	if account.BusinessEmail == EMPTY || account.PhoneNumber == EMPTY || account.BusinessName == EMPTY {
		return service_errors.ErrMisconfiguredAccountParameters
	}

	return nil
}

// ValidateAccountNotNil ensures the account object is not nil
func (db *Db) ValidateAccountNotNil(ctx context.Context, account *models.MerchantAccountORM) error {
	if account == nil {
		return service_errors.ErrInvalidAccount
	}

	return nil
}

// ValidateAccountIds validates the existence of various ids associated with the account
func (db *Db) ValidateAccountIds(ctx context.Context, account *models.MerchantAccountORM) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err.Error())
		return err
	}

	if account.AuthnAccountId == 0 || account.StripeConnectedAccountId == EMPTY || account.EmployerId == 0 {
		return service_errors.ErrMisconfiguredIds
	}

	return nil
}

func (db *Db) AccountActive(account *models.MerchantAccountORM) bool {
	if account == nil || !account.IsActive {
		return false
	}

	return true
}

// UpdateAccountOnboardStatus updates the onboarding status of a merchant account
func (db *Db) UpdateAccountOnboardStatus(ctx context.Context, account *models.MerchantAccountORM) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err.Error())
		return err
	}

	switch account.AccountOnboardingDetails {
	case int32(models.OnboardingStatus_OnboardingNotStarted):
		db.feelguudsOnboardingStarted(account)
	case int32(models.OnboardingStatus_FeelGuudOnboardingStarted):
		db.feelguudsOnboardingCompleted(account)
	case int32(models.OnboardingStatus_StripeOnboardingStarted):
		db.stripeOnboardingComplete(account)
	case int32(models.OnboardingStatus_BCorpOnboardingStarted):
		db.bcorpOnboardingComplete(account)
	case int32(models.OnboardingStatus_CatalogueOnboardingStarted):
		db.catalogueOnboardingStarted(account)
	default:
		db.onboardingNotStarted(account)
	}

	return nil
}

// onboardingNotStarted sets the account onboarding state to onboarding not started
func (db *Db) onboardingNotStarted(account *models.MerchantAccountORM) {
	account.AccountOnboardingDetails = int32(models.OnboardingStatus_OnboardingNotStarted)
	account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
}

// catalogueOnboardingStarted sets the account onboarding state to catalogue onboarding completed
func (db *Db) catalogueOnboardingStarted(account *models.MerchantAccountORM) {
	account.AccountOnboardingDetails = int32(models.OnboardingStatus_CatalogueOnboardingCompleted)
	account.AccountOnboardingState = int32(models.MerchantAccountState_ActiveAndOnboarded)
}

// bcorpOnboardingComplete sets the account onboarding state to b-corp onboarding completed
func (db *Db) bcorpOnboardingComplete(account *models.MerchantAccountORM) {
	account.AccountOnboardingDetails = int32(models.OnboardingStatus_BCorpOnboardingCompleted)
	account.AccountOnboardingState = int32(models.MerchantAccountState_ActiveAndOnboarded)
}

// stripeOnboardingComplete sets the account onboarding state to stripe onboarding completed
func (db *Db) stripeOnboardingComplete(account *models.MerchantAccountORM) {
	account.AccountOnboardingDetails = int32(models.OnboardingStatus_StripeOnboardingCompleted)
	account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
}

// feelguudsOnboardingCompleted sets the account onboarding state to feelguuds onboarding completed
func (db *Db) feelguudsOnboardingCompleted(account *models.MerchantAccountORM) {
	account.AccountOnboardingDetails = int32(models.OnboardingStatus_FeelGuudOnboardingCompleted)
	account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
}

// feelguudsOnboardingCompleted sets the account onboarding state to feelguuds onboarding started
func (db *Db) feelguudsOnboardingStarted(account *models.MerchantAccountORM) {
	account.AccountOnboardingDetails = int32(models.OnboardingStatus_FeelGuudOnboardingStarted)
	account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
}
