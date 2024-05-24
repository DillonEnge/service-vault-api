package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/authz"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var (
	authzBypass = []string{
		"POST /request/access",
		"POST /request/code",
		"POST /report",
		"GET /services",
	}
	totalBypass = []string{
		"GET /health",
		"GET /signin",
	}
)

func Authorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		if slices.Contains(totalBypass, route) {
			slog.Info("Backdoor path detected, bypassing authn and authz check")

			next.ServeHTTP(w, r)

			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Error("Failed to find Authorization header")
			errors.Unauthorized(w, r)

			return
		}

		if slices.Contains(authzBypass, route) {
			slog.Info("Backdoor path detected, bypassing authz check")
			next.ServeHTTP(w, r)

			return
		}

		ctx := r.Context()

		cdoorClient := ctx.Value(consts.CtxKeyCasdoorClient).(*casdoorsdk.Client)

		claims, err := cdoorClient.ParseJwtToken(strings.Split(authHeader, " ")[1])
		if err != nil {
			errors.InternalServerError(w, r, err)

			return
		}

		cbinClient := ctx.Value(consts.CtxKeyCasbinClient).(*authz.Client)

		cbinClient.Enforcer.LoadPolicy()

		if !cbinClient.GetAuthorization(claims.Email, r.URL.Path, r.Method) {
			errors.Forbidden(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}
