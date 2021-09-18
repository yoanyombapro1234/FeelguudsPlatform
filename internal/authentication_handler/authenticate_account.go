package authentication_handler

import (
	"context"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
)

type AuthenticateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateAccountResponse struct {
	Token        string `json:"token"`
	Code         int    `json:"code"`
	ErrorMessage string `json:"message"`
}

// AuthenticateAccountHandler godoc
// @Summary Account Authentication
// @Description authenticates a user account based on provided credentials against the authentication service
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (c *AuthenticationComponent) AuthenticateAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), c.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing
	var (
		authenticationAccountReq AuthenticateAccountRequest
	)

	err := helper.DecodeJSONBody(w, r, &authenticationAccountReq)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if authenticationAccountReq.Email == EMPTY || authenticationAccountReq.Password == EMPTY {
		helper.ErrorResponse(w, "invalid email or password. please provide valid input parameters", http.StatusBadRequest)
		return
	}

	email, password := authenticationAccountReq.Email, authenticationAccountReq.Password
	// invoke authentication service
	token, err := c.AuthenticateAccount(ctx, email, password)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, &AuthenticateAccountResponse{Token: token})
}

// AuthenticateAccount attempts to authenticate a user based on provided credentials
func (c *AuthenticationComponent) AuthenticateAccount(ctx context.Context, email, password string) (string, error) {
	if _, err := c.IsPasswordOrEmailInValid(email, password); err != nil {
		return "", err
	}

	token, err := c.Client.LoginAccount(email, password)
	if err != nil {
		c.Logger.Error(err.Error())
		return "", err
	}

	if _, err := c.checkJwtTokenForInValidity(token); err != nil {
		c.Logger.Error(err.Error())
	}

	c.Logger.Info("successfully authenticated user account")
	return token, nil
}
