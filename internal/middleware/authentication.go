package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	"go.uber.org/zap"
)

type contextKey struct {
	name string
}

type AuthnMW struct {
	client *core_auth_sdk.Client
	logger *zap.Logger
}

var ctxKey *contextKey

// NewAuthnMiddleware returns a new instance of the authentication middleware
func NewAuthnMiddleware(c *core_auth_sdk.Client, log *zap.Logger, serviceName string) *AuthnMW {
	if serviceName == "" {
		panic("service name should be provided")
	}

	ctxKey = &contextKey{serviceName}
	return &AuthnMW{client: c, logger: log}
}

// AuthenticationMiddleware wraps the authentication middleware around an http call
func (mw *AuthnMW) AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()
		authorization := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")
		decodedToken, err := mw.client.SubjectFrom(token)
		if err != nil {
			mw.logger.Error(fmt.Sprintf("user not authenticated. error: %s", err.Error()))
			return
		}

		ctx = context.WithValue(ctx, ctxKey, decodedToken)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// IsAuthenticated returns wether or not the user is authenticated.
// REQUIRES Middleware to have run.
func IsAuthenticated(ctx context.Context) bool {
	return ctx.Value(ctxKey) != nil
}

// GetTokenFromCtx extracts the token from the context
func GetTokenFromCtx(ctx context.Context) (string, error) {
	if IsAuthenticated(ctx) {
		token, ok := ctx.Value(ctxKey).(string)
		if !ok {
			return "", errors.New("token cannot be converted to string")
		}

		return token, nil
	}

	return "", errors.New("token not found in context")
}

// InjectContextWithMockToken injects a token into the context. Useful for testing
func InjectContextWithMockToken(ctx context.Context, token string, serviceName string) context.Context {
	ctxKey = &contextKey{serviceName}
	return context.WithValue(ctx, ctxKey, token)
}
