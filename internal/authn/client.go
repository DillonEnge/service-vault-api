package authn

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	ierrors "github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var (
	ErrAuthHeaderNotFound = errors.New("Auth header not found")
)

func GetClaims(w http.ResponseWriter, r *http.Request) (*casdoorsdk.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		slog.Error("Failed to find Authorization header")
		ierrors.Unauthorized(w, r)

		return nil, ErrAuthHeaderNotFound
	}

	ctx := r.Context()

	cdoorClient := ctx.Value(consts.CtxKeyCasdoorClient).(*casdoorsdk.Client)

	claims, err := cdoorClient.ParseJwtToken(strings.Split(authHeader, " ")[1])
	if err != nil {
		ierrors.InternalServerError(w, r, err)

		return nil, err
	}

	return claims, nil
}
