package middleware

import (
	"math/rand"
	"net/http"
	"time"
)

func RandomErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().Unix())
		if rand.Int31n(3) == 0 {

			errors := []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusConflict}
			w.WriteHeader(errors[rand.Intn(len(errors))])
			return
		}
		next.ServeHTTP(w, r)
	})
}
