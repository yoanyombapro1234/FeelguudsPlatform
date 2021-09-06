package authentication_handler

import (
	"context"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
)

type LogoutAccountResponse struct {
	Code  int    `json:"code"`
	ErrorMessage string `json:"message"`
}

// LogoutAccountHandler godoc
// @Summary Log out of Account
// @Description logs user account out of the system from the context of the authentication service
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (c *AuthenticationComponent) LogoutAccountHandler(w http.ResponseWriter, r *http.Request) {
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
	if err = c.LogoutAccount(ctx, Id); err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, LogoutAccountResponse{})
}

func (c *AuthenticationComponent) LogoutAccount(ctx context.Context, Id uint32) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	// TODO: think about how to handle this failed call
	if err := c.Client.LogOutAccount(); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	return nil
}
