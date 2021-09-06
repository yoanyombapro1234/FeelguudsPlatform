package authentication_handler

import (
	"time"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-logging"
	"go.uber.org/zap"
)

const username string = "feelguuds"
const password string = "feelguuds"
const audience string = "localhost"
const issuer string = "http://localhost:8000"
const origin string = "http://localhost"
const privateBaseUrl string = "http://localhost:8000"
const svcName = "mock"

// InitializeMockAuthenticationComponent initializes a mock authentication component assuming certain pre-conditions are met
func InitializeMockAuthenticationComponent() *AuthenticationComponent {
	logInstance := core_logging.New("info")
	defer logInstance.ConfigureLogger()
	log := logInstance.Logger

	return mockAuthenticationComponent(log, svcName)
}

// MockAuthenticationComponent configures a mock authentication component
func mockAuthenticationComponent(log *zap.Logger, serviceName string) *AuthenticationComponent {
	httpTimeout := 300 * time.Millisecond

	return NewAuthenticationComponent(&AuthenticationParams{
		AuthConfig: &core_auth_sdk.Config{
			Issuer:         issuer,
			PrivateBaseURL: privateBaseUrl,
			Audience:       audience,
			Username:       username,
			Password:       password,
			KeychainTTL:    0,
		},
		AuthConnectionConfig: &core_auth_sdk.RetryConfig{
			MaxRetries:       1,
			MinRetryWaitTime: 200 * time.Millisecond,
			MaxRetryWaitTime: 300 * time.Millisecond,
			RequestTimeout:   500 * time.Millisecond,
		},
		Logger: log,
		Origin: origin,
	}, serviceName, httpTimeout)
}
