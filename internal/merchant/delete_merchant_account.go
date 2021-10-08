package merchant

import (
	"context"
	"errors"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
)

type DeleteMerchantAccountResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"message"`
}

// DeleteMerchantAccountHandler godoc
// @Summary Deletes a merchant account
// @Description coordinates interactions across multiple services to delete a merchant account
// @Tags HTTP API
// @Produce html
// @Router / [delete]
// @Success 200 {string} string "OK"
func (m *MerchantAccountComponent) DeleteMerchantAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	id, err := helper.ExtractIDFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.DeleteMerchantAccount(ctx, id); err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, DeleteMerchantAccountResponse{})
}

func (m MerchantAccountComponent) DeleteMerchantAccount(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("merchant account id cannot be 0")
	}

	if _, err := m.Db.DeactivateMerchantAccount(ctx, id); err != nil {
		return err
	}

	return nil
}
