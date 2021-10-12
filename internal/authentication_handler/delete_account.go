package authentication_handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"go.uber.org/zap"
)

type DeleteAccountResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"message"`
}

// DeleteAccountHandler godoc
// @Summary Delete Account
// @Description deletes user account in the authentication service
// @Tags HTTP API
// @Produce html
// @Router / [delete]
// @Success 200 {string} string "OK"
func (c *AuthenticationComponent) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), c.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing and logs

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
	if err := c.DeleteAccount(ctx, uint32(Id)); err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, DeleteAccountResponse{})
}

// DeleteAccount attempts to archive an account from the context of the authentication service (authn)
func (c *AuthenticationComponent) DeleteAccount(ctx context.Context, Id uint32) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	accountId := strconv.Itoa(int(Id))
	if err := c.Client.ArchiveAccount(accountId); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	c.Logger.Info("Successfully deleted user account", zap.Int("Id", int(Id)))
	return nil
}
