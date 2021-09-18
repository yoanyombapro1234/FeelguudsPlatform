package authentication_handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"go.uber.org/zap"
)

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

type UpdateEmailResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"message"`
}

// UpdateEmailHandler godoc
// @Summary Update Account Email
// @Description updates the email account from the context of the authentication service
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (c *AuthenticationComponent) UpdateEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), c.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing and logs
	var (
		updateEmailReq UpdateEmailRequest
	)

	err := helper.DecodeJSONBody(w, r, &updateEmailReq)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateEmailReq.Email == EMPTY {
		helper.ErrorResponse(w, "invalid email. please provide valid input parameters", http.StatusBadRequest)
		return
	}

	email := updateEmailReq.Email

	// get Id from authentication service and validate
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
	if err = c.UpdateAccount(ctx, Id, email); err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, UpdateEmailResponse{})
}

// UpdateAccount updates a user account's credentials
func (c *AuthenticationComponent) UpdateAccount(ctx context.Context, Id uint32, email string) error {
	if err, _ := c.isValidID(Id); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	if err, _ := c.isValidEmail(email); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	accountId := strconv.Itoa(int(Id))
	if err := c.Client.Update(accountId, email); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	c.Logger.Info("Successfully updated user account", zap.Int("Id", int(Id)))
	return nil
}
