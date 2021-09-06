package authentication_handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"go.uber.org/zap"
)

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountResponse struct {
	Id    uint32 `json:"id"`
	Code  int    `json:"code"`
	ErrorMessage string `json:"message"`
}

// CreateAccountHandler godoc
// @Summary Create Account
// @Description creates a new user account in the authentication service
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (c *AuthenticationComponent) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), c.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing
	var (
		createAccountReq CreateAccountRequest
	)

	err := helper.DecodeJSONBody(w, r, &createAccountReq)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if createAccountReq.Email == EMPTY || createAccountReq.Password == EMPTY {
		helper.ErrorResponse(w, "invalid email or password. Please provide the proper credentials", http.StatusBadRequest)
		return
	}

	email, password := createAccountReq.Email, createAccountReq.Password
	// invoke authentication service
	id, err := c.CreateAccount(ctx, email, password, false)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, &CreateAccountResponse{Id: id})
}

// CreateAccount attempts to create a user account against the authentication service
func (c *AuthenticationComponent) CreateAccount(ctx context.Context, email, password string, accountLocked bool) (uint32, error) {
	if _, err := c.IsPasswordOrEmailInValid(email, password); err != nil {
		return 0, err
	}

	accountId, err := c.Client.ImportAccount(email, password, accountLocked)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("failed to create account against authentication service for user %s. error: %s", email, err.Error()))
		return 0, err
	}

	c.Logger.Info("Successfully created user account", zap.Int("Id", int(accountId)))
	return uint32(accountId), nil
}
