package authentication_handler

import (
	"context"
	"net/http"
	"strconv"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"go.uber.org/zap"
)

type GetAccountResponse struct {
	Account *core_auth_sdk.Account `json:"account"`
	Code  int    `json:"code"`
	ErrorMessage string `json:"message"`
}

// GetAccountHandler godoc
// @Summary Get Account
// @Description gets a user account from the context of the authentication service
// @Tags HTTP API
// @Produce html
// @Router / [get]
// @Success 200 {string} string "OK"
func (c *AuthenticationComponent) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), c.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing and logs
	var (
		account *core_auth_sdk.Account
	)

	Id, err := helper.ExtractIDFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if Id == 0 {
		helper.ErrorResponse(w, "invalid user account id. please provide valid input parameters", http.StatusBadRequest)
		return
	}

	// invoke authentication service
	account, err = c.GetAccount(ctx, Id)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, GetAccountResponse{Account: account})
}


// GetAccount obtains a user account from the context of the authentications service (authn) based on a provided user id
func (c *AuthenticationComponent) GetAccount(ctx context.Context, Id uint32) (*core_auth_sdk.Account, error) {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return nil, err
	}

	accountId := strconv.Itoa(int(Id))
	account, err := c.Client.GetAccount(accountId)
	if err != nil {
		c.Logger.Error(err.Error())
		return nil, err
	}

	c.Logger.Info("Successfully obtained user account", zap.Int("Id", int(Id)))
	return account, nil
}
