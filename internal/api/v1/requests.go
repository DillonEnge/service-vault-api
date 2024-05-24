package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/authn"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/DillonEnge/service-vault-api/internal/mail"
	"github.com/DillonEnge/service-vault-api/internal/sql/requests"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CreateCodeRequestParams struct {
	ServiceName pgtype.Varchar `json:"service_name"`
}

func GetAllRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pool := ctx.Value(consts.CtxKeyDbPool).(*pgxpool.Pool)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to aquire conn from pool")
		errors.InternalServerError(w, r, err)
		return
	}
	defer conn.Release()

	q := requests.NewQuerier(conn.Conn())

	slog.Info("performing query")

	res, err := q.GetAllRequests(ctx)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	slog.Info("generating glob")

	glob, err := json.Marshal(res)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(glob)
}

func PostCodeRequest(w http.ResponseWriter, r *http.Request) {
	var p CreateCodeRequestParams

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	c, err := authn.GetClaims(w, r)
	if err != nil {
		slog.Error("Failed to get claims")
		errors.InternalServerError(w, r, err)
		return
	}

	var email pgtype.Varchar

	err = email.Set(c.Email)
	if err != nil {
		slog.Error("Failed to convert email str to varchar")
		errors.InternalServerError(w, r, err)
		return
	}

	pool := ctx.Value(consts.CtxKeyDbPool).(*pgxpool.Pool)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to aquire conn from pool")
		errors.InternalServerError(w, r, err)
		return
	}
	defer conn.Release()

	q := requests.NewQuerier(conn.Conn())

	_, err = q.CreateCodeRequest(ctx, p.ServiceName, email)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	data := map[string]string{"email": c.Email, "service": p.ServiceName.String}

	recipients := []string{"dillon.enge@gmail.com"}

	subject := "Vault 2FA Access Code Request"

	err = mail.SendMail("tmpl/request_code.html", data, recipients, subject)
	if err != nil {
		slog.Error("Failed to send mail")
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func PostAccessRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := authn.GetClaims(w, r)
	if err != nil {
		slog.Error("Failed to get claims")
		return
	}

	var email pgtype.Varchar

	err = email.Set(c.Email)
	if err != nil {
		slog.Error("Failed to convert email str to varchar")
		errors.InternalServerError(w, r, err)
		return
	}

	pool := ctx.Value(consts.CtxKeyDbPool).(*pgxpool.Pool)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to aquire conn from pool")
		errors.InternalServerError(w, r, err)
		return
	}
	defer conn.Release()

	q := requests.NewQuerier(conn.Conn())

	_, err = q.CreateAccessRequest(ctx, email)
	if err != nil {
		errors.InternalServerError(w, r, err)
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

func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Error("Failed to find mandatory query param 'id'")
		errors.InternalServerError(w, r, nil)
		return
	}

	var uuid pgtype.UUID

	err := uuid.Set(id)
	if err != nil {
		slog.Error("Failed to convert id to uuid")
		errors.InternalServerError(w, r, err)
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

	q := requests.NewQuerier(conn.Conn())

	_, err = q.DeleteRequest(ctx, uuid)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
