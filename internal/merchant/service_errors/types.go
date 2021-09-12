package service_errors

import (
	"errors"
)

// List of values that ErrorCode can take.
const (
	ErrorCodeAccountAlreadyExists           ErrorCode = "account_already_exists"
	ErrorCodeMisconfiguredMerchantAccountID ErrorCode = "misconfigured_merchant_account_id"
	ErrorCodeAccountInvalid                 ErrorCode = "account_invalid"
	ErrorCodeMisconfiguredAccountParameters ErrorCode = "misconfigured_merchant_account_parameters"
	ErrorCodeDatabaseConnectionFailure      ErrorCode = "database_connection_failure"
	ErrorCodeDatabaseMigrationFailure       ErrorCode = "database_migration_failure"
	ErrorCodeInvalidInput                   ErrorCode = "invalid_input"
	ErrorCodeInvalidEnvironmentVariables    ErrorCode = "invalid_environment_variables"
	ErrorCodeAuthenticationRequired         ErrorCode = "authentication_required"
)

var (
	ErrMisconfiguredAccountParameters           = errors.New("misconfigured merchant account parameters")
	ErrMisconfiguredIds                         = errors.New("misconfigured merchant account id")
	ErrFailedToConnectToDatabase                = errors.New("failed to connect to database")
	ErrFailedToPerformDatabaseMigrations        = errors.New("failed to perform database migrations")
	ErrInvalidInputArguments                    = errors.New("invalid input arguments")
	ErrInvalidEnvironmentVariableConfigurations = errors.New("invalid environment variable configurations")
	ErrFailedToStartGRPCServer                  = errors.New("failed to start grpc server")
	ErrHttpServerFailedGracefuleShutdown        = errors.New("http server failed to perform graceful shutdown")
	ErrHttpsServerFailedGracefuleShutdown       = errors.New("https server failed to perform graceful shutdown")
	ErrHttpServerCrashed                        = errors.New("Http Server crashed")
	ErrHttpsServerCrashed                       = errors.New("Https Server crashed")
	ErrSwaggerGenError                          = errors.New("swagger generation error")
	ErrFailedToWatchConfigDirectory             = errors.New("failed to watch config directory")
	ErrExceededMaxRetryAttempts                 = errors.New("exceeded max retry attemps")
	ErrInvalidAccount                           = errors.New("invalid account. account contains invalid fields")
	ErrFailedToReactivateAccount                = errors.New("failed to reactivate existing account")
	ErrDistributedTransactionError              = errors.New("distributed transaction error. failed to successfully perform a distributed operations")
	ErrFailedToUpdateAccountActiveStatus        = errors.New("failed to updated account active status")
	ErrAccountDoesNotExist                      = errors.New("account does not exist")
	ErrAccountAlreadyExist                      = errors.New("account already exists")
	ErrFailedToConvertFromOrmType               = errors.New("failed to perform conversion from Orm type")
	ErrFailedToConvertToOrmType                 = errors.New("failed to perform conversion to Orm type")
	ErrFailedToConfigureSaga                    = errors.New("failed to configure saga")
	ErrSagaFailedToExecuteSuccessfully          = errors.New("saga failed to execute successfully")
	ErrFailedToHashPassword                     = errors.New("failed to hash password")
	ErrFailedToCreateAccount                    = errors.New("failed to create account")
	ErrFailedToDeleteBusinessAccount            = errors.New("failed to delete business account")
	ErrFailedToUpdateAccountEmail               = errors.New("failed to updated business account email through distributed transaction")
	ErrFailedToSaveUpdatedAccountRecord         = errors.New("failed to save updated business account record")
	ErrCannotUpdatePassword                     = errors.New("cannot update password field")
	ErrCannotConfigureAccount                   = errors.New("cannot configure account")
	ErrFailedToCastToType                       = errors.New("failed to cast to tx result to type after tx completion")
	ErrUnableToObtainBusinessAccounts           = errors.New("unable to obtain business accounts")
	ErrUnauthorizedRequest                      = errors.New("unauthorized request")
	ErrFailedToCreateAccountInAuthHdlrSvc       = errors.New("failed to create an account in authentication handler service")
)

func NewError(msg string) error {
	return errors.New(msg)
}
