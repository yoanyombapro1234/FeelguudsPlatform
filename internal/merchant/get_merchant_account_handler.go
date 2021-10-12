package merchant

import (
	"context"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/models"
)

type GetMerchantAccountResponse struct {
	MerchantAccount *models.MerchantAccount `json:"merchant_account"`
}

// GetMerchantAccountHandler godoc
// @Summary Gets a merchant account if it exists
// @Description returns a merchant account if it exists
// @Tags HTTP API
// @Produce html
// @Router / [delete]
// @Success 200 {string} string "OK"
func (m *AccountComponent) GetMerchantAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	id, err := helper.ExtractIDFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	acct, err := m.Db.GetMerchantAccountById(ctx, id, true)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.JSONResponse(w, GetMerchantAccountResponse{acct})
}
