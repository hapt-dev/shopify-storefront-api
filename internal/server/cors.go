package server

import (
	"net/http"
)

var (
	corsAllowedMethods = "OPTIONS, HEAD, GET, PUT, PATCH, POST, DELETE"
	corsAllowedHeaders = "Origin, X-Requested-With, Content-Type, Accept, Authorization, authentication"
)

// cors matches NestJS payment-api enableCors settings.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", corsAllowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", corsAllowedHeaders)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
