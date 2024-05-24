package v1

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/authn"
	"github.com/DillonEnge/service-vault-api/internal/authz"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/DillonEnge/service-vault-api/internal/mail"
)

func GetAccess(w http.ResponseWriter, r *http.Request) {
	c, err := authn.GetClaims(w, r)
	if err != nil {
		slog.Error("Failed to get claims")
		return
	}

	data := map[string]string{"email": c.Email}

	recipients := []string{"dillon.enge@gmail.com"}

	subject := "Vault Access Request"

	err = mail.SendMail("tmpl/request.html", data, recipients, subject)
	if err != nil {
		slog.Error("Failed to send email")
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func PostAccess(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		slog.Error("Failed to find mandatory 'email' param in query string")
		errors.InternalServerError(w, r, nil)
		return
	}

	shouldNotify := false

	notify := r.URL.Query().Get("notify")
	if notify == "true" {
		shouldNotify = true
	}

	cbinClient := r.Context().Value(consts.CtxKeyCasbinClient).(*authz.Client)

	ok, err := cbinClient.Enforcer.AddPolicy(email, "/service", "GET")
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	if !ok {
		slog.Error("Policy rule already present")
		errors.InternalServerError(w, r, nil)
		return
	}

	err = cbinClient.Enforcer.SavePolicy()
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	if !shouldNotify {
		w.WriteHeader(http.StatusOK)
		return
	}

	email = strings.ToLower(email)

	data := map[string]string{"email": email}

	recipients := []string{email}

	subject := "Vault Access Granted"

	err = mail.SendMail("tmpl/grant.html", data, recipients, subject)
	if err != nil {
		slog.Error("Failed to send email")
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
