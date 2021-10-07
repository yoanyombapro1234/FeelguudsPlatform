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

	if account.AuthnAccountId == 0 || account.StripeAccountId == 0 || account.StripeConnectedAccountId == EMPTY || account.EmployerId == 0 {
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
	// not started onboarding
	case int32(models.OnboardingStatus_OnboardingNotStarted):
		account.AccountOnboardingDetails = int32(models.OnboardingStatus_FeelGuudOnboarding)
		account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
		// completed onboarding with feelguud
	case int32(models.OnboardingStatus_FeelGuudOnboarding):
		account.AccountOnboardingDetails = int32(models.OnboardingStatus_StripeOnboarding)
		account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
		// completed onboarding with stripe
	case int32(models.OnboardingStatus_StripeOnboarding):
		account.AccountOnboardingDetails = int32(models.OnboardingStatus_CatalogueOnboarding)
		account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
		// completed onboarding catalogue
	case int32(models.OnboardingStatus_CatalogueOnboarding):
		account.AccountOnboardingDetails = int32(models.OnboardingStatus_BCorpOnboarding)
		account.AccountOnboardingState = int32(models.MerchantAccountState_ActiveAndOnboarded)
	default:
		account.AccountOnboardingDetails = int32(models.OnboardingStatus_OnboardingNotStarted)
		account.AccountOnboardingState = int32(models.MerchantAccountState_PendingOnboardingCompletion)
	}

	return nil
}
