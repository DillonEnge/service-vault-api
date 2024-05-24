package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/authn"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/DillonEnge/service-vault-api/internal/mail"
	"github.com/DillonEnge/service-vault-api/internal/sql/services"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostReportParams struct {
	ServiceName pgtype.Varchar `json:"service_name"`
}

func ReportService(w http.ResponseWriter, r *http.Request) {
	var p PostReportParams
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		slog.Error("Failed to decode json body")
		errors.InternalServerError(w, r, err)
		return
	}

	c, err := authn.GetClaims(w, r)
	if err != nil {
		slog.Error("Failed to get claims")
		return
	}

	ctx := r.Context()

	pool := ctx.Value(consts.CtxKeyDbPool).(*pgxpool.Pool)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to aquire conn from pool")
		errors.InternalServerError(w, r, err)
		return
	}
	defer conn.Release()

	q := services.NewQuerier(conn.Conn())

	_, err = q.ReportService(ctx, p.ServiceName.String)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	data := map[string]string{"email": c.Email, "service": p.ServiceName.String}

	recipients := []string{"dillon.enge@gmail.com"}

	subject := "Vault Stale Password Report"

	err = mail.SendMail("tmpl/report.html", data, recipients, subject)
	if err != nil {
		slog.Error("Failed to send email")
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ResetReportedService(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	if service == "" {
		slog.Error("Failed to find mandatory query param 'service'")
		errors.InternalServerError(w, r, nil)
		return
	}

	ctx := r.Context()

	pool := ctx.Value(consts.CtxKeyDbPool).(*pgxpool.Pool)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to aquire conn from pool")
		errors.InternalServerError(w, r, err)
		return
	}
	defer conn.Release()

	q := services.NewQuerier(conn.Conn())

	_, err = q.ResetReportedService(ctx, service)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
