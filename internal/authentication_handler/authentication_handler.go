package authentication_handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/giantswarm/retry-go"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	"go.uber.org/zap"
)

type AuthenticationServiceInterface interface {
	AuthenticateAccount(ctx context.Context, email, password string) (string, error)
	CreateAccount(ctx context.Context, email, password string, accountLocked bool) (uint32, error)
	DeleteAccount(ctx context.Context, Id uint32) error
	GetAccount(ctx context.Context, Id uint32) (*core_auth_sdk.Account, error)
	LockAccount(ctx context.Context, Id uint32) error
	UnLockAccount(ctx context.Context, Id uint32) error
	UpdateAccount(ctx context.Context, Id uint32, email string) error
	LogoutAccount(ctx context.Context, Id uint32) error
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
}

var _ AuthenticationServiceInterface = (*AuthenticationComponent)(nil)

// NewAuthenticationComponent returns an authentication component to the caller
func NewAuthenticationComponent(params *AuthenticationParams, serviceName string) *AuthenticationComponent {
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
	if err = ConnectToDownstreamService(logger, authnClient, response); err != nil {
		logger.Fatal(fmt.Sprintf("failed to actually connect to authn client. error - %s", err.Error()))
	}

	serviceMetrics := NewServiceMetrics(serviceName)

	return &AuthenticationComponent{Client: authnClient, Logger: params.Logger, Metric: serviceMetrics}
}

// ConnectToDownstreamService attempts to connect to a downstream service
func ConnectToDownstreamService(logger *zap.Logger, client *core_auth_sdk.Client, response chan interface{}) error {
	return retry.Do(
		func(conn chan<- interface{}) func() error {
			return func() error {
				data, err := client.ServerStats()
				if err != nil {
					logger.Error("failed to connect to authentication service", zap.Error(err))
					return err
				}

				logger.Info("data", zap.Any("result", data))

				response <- data
				return nil
			}
		}(response),
		retry.MaxTries(5),
		retry.Timeout(time.Millisecond*time.Duration(10)),
		retry.Sleep(time.Millisecond*time.Duration(10)))
}
