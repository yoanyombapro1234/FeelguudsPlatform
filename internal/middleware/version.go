package middleware

import (
	"net/http"

	"github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version"
)

func VersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-API-Version", version.VERSION)
		r.Header.Set("X-API-Revision", version.REVISION)

		next.ServeHTTP(w, r)
	})
}
