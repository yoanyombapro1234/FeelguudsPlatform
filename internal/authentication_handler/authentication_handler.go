package authentication_handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/giantswarm/retry-go"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	"go.uber.org/zap"
)

// AuthenticationServiceInterface provides an interface definition specific to authentication
type AuthenticationServiceInterface interface {
	AuthenticateAccount(ctx context.Context, email, password string) (string, error)
	AuthenticateAccountHandler(w http.ResponseWriter, r *http.Request)
	CreateAccount(ctx context.Context, email, password string, accountLocked bool) (uint32, error)
	CreateAccountHandler(w http.ResponseWriter, r *http.Request)
	DeleteAccount(ctx context.Context, Id uint32) error
	DeleteAccountHandler(w http.ResponseWriter, r *http.Request)
	GetAccount(ctx context.Context, Id uint32) (*core_auth_sdk.Account, error)
	GetAccountHandler(w http.ResponseWriter, r *http.Request)
	LockAccount(ctx context.Context, Id uint32) error
	LockAccountHandler(w http.ResponseWriter, r *http.Request)
	UnLockAccount(ctx context.Context, Id uint32) error
	UnLockAccountHandler(w http.ResponseWriter, r *http.Request)
	UpdateAccount(ctx context.Context, Id uint32, email string) error
	UpdateEmailHandler(w http.ResponseWriter, r *http.Request)
	LogoutAccount(ctx context.Context, Id uint32) error
	LogoutAccountHandler(w http.ResponseWriter, r *http.Request)
}

// AuthenticationParams encompases the required entries necessary to configure a client connection to the authn service
type AuthenticationParams struct {
	// AuthConfig is comprised with security parameters necessary for connecting to the authn service
	AuthConfig *core_auth_sdk.Config
	// AuthnConnectionConfig defines the various retry configurations that will dictate the retry logic to engage in in the face of an http failure
	AuthConnectionConfig *core_auth_sdk.RetryConfig
	// Logger is the logger object
	Logger *zap.Logger
	// Origin is the origin server from which requests originate from
	Origin string
}

// AuthenticationComponent provides a wrapper around the authn client for more robust configurability per our use cases
type AuthenticationComponent struct {
	// Client is a connection handler to the authn service
	Client *core_auth_sdk.Client
	// Logger is the logger object used by this component
	Logger *zap.Logger
	// Metric specific to this module
	Metric *ServiceMetrics
	// Duration of any expected http call
	HttpTimeout time.Duration
}

var _ AuthenticationServiceInterface = (*AuthenticationComponent)(nil)

// NewAuthenticationComponent returns an authentication component to the caller
func NewAuthenticationComponent(params *AuthenticationParams, serviceName string, httpRequestTimeout time.Duration) *AuthenticationComponent {
	if params == nil {
		log.Fatal(ErrInvalidInputArguments.Error())
	}

	if params.Origin == "" || params.AuthConfig == nil || params.AuthConnectionConfig == nil || params.Logger == nil {
		log.Fatal(ErrInvalidInputArguments.Error())
	}

	logger := params.Logger
	authnClient, err := core_auth_sdk.NewClient(*params.AuthConfig, params.Origin, params.AuthConnectionConfig)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s, error - %s", "unable to configure connection to authn client", err.Error()))
	}

	var response = make(chan interface{}, 1)
	if err = ConnectToAuthService(logger, authnClient, response); err != nil {
		logger.Fatal(fmt.Sprintf("failed to actually connect to authn client. error - %s", err.Error()))
	}

	logger.Info("successfully connected to authentication service")

	serviceMetrics := NewServiceMetrics(serviceName)

	return &AuthenticationComponent{Client: authnClient, Logger: params.Logger, Metric: serviceMetrics, HttpTimeout: httpRequestTimeout}
}

// ConnectToAuthService attempts to connect to a downstream service
func ConnectToAuthService(logger *zap.Logger, client *core_auth_sdk.Client, response chan interface{}) error {
	return retry.Do(
		func(conn chan<- interface{}) func() error {
			return func() error {
				res, err := client.ServerStats()
				if err != nil {
					logger.Error("failed to connect to authentication service", zap.Error(err))
					return err
				}

				body, err := io.ReadAll(res.Body)
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						logger.Error(err.Error())
					}
				}(res.Body)

				response <- body
				return nil
			}
		}(response),
		retry.MaxTries(5),
		retry.Timeout(time.Millisecond*time.Duration(10)),
		retry.Sleep(time.Millisecond*time.Duration(10)))
}
