package middleware

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin, ok := os.LookupEnv("ALLOW_ORIGIN")
		if !ok {
			slog.Error("Failed to load CASDOOR_ENDPOINT from env")
			errors.InternalServerError(w, r, nil)
			return
		}

		headers := w.Header()
		headers.Add("Access-Control-Allow-Origin", origin)
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Add("Access-Control-Allow-Headers", "Content-Type, Accept, Origin, Authorization")
		headers.Add("Access-Control-Allow-Methods", "GET,POST,PATCH,OPTIONS,DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)

			return
		}

		next.ServeHTTP(w, r)
	})
}
