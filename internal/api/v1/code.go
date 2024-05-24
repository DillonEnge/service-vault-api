package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/mail"
)

type PostCodeRequestParams struct {
	Email       string `json:"email"`
	ServiceName string `json:"service_name"`
	Code        string `json:"code"`
}

func PostCode(w http.ResponseWriter, r *http.Request) {
	var p PostCodeRequestParams

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	data := map[string]string{
		"email":   p.Email,
		"service": p.ServiceName,
		"code":    p.Code,
	}

	recipients := []string{p.Email}

	subject := "Vault 2FA Access Code Granted"

	err = mail.SendMail("tmpl/grant_code.html", data, recipients, subject)
	if err != nil {
		slog.Error("Failed to send mail")
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
