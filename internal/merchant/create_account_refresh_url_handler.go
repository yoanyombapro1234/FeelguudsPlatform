package merchant

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
)

// CreateAccountRefreshUrlHandler godoc
// @Summary Serves as the handler for the refresh url used throughout the merchant account onboarding process
// @Description resets the onboarding process for a given user
// @Tags HTTP API
// @Produce html
// @Router / [post]
// @Success 200 {string} string "OK"
func (m *MerchantAccountComponent) CreateAccountRefreshUrlHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.HttpTimeout)
	defer cancel()

	// TODO: emit metrics and add distributed tracing
	connectedAcctId, err := helper.ExtractStripeConnectedAccountIdFromRequest(r)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshUrl := fmt.Sprintf("%s/%s", m.BaseRefreshUrl, connectedAcctId)
	returnUrl := fmt.Sprintf("%s/%s", m.BaseReturnUrl, connectedAcctId)

	// pull the merchant account from the database
	acct, err := m.Db.FindMerchantAccountByStripeAccountId(ctx, connectedAcctId)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	stripeRedirectUri, err := m.getStripeRedirectURI(acct, connectedAcctId, refreshUrl, returnUrl)
	if err != nil {
		helper.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, stripeRedirectUri, http.StatusOK)
}
