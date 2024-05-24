package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/DillonEnge/service-vault-api/internal/sql/services"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PatchServicePasswordBody struct {
	Password    pgtype.Varchar `json:"password"`
	ServiceName string         `json:"service_name"`
}

func GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		slog.Error("Failed to find mandatory 'name' param in query string")
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

	res, err := q.GetService(ctx, name)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	glob, err := json.Marshal(res)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(glob)
}

func GetAllServices(w http.ResponseWriter, r *http.Request) {
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

	res, err := q.GetAllServices(ctx)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	glob, err := json.Marshal(res)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(glob)
}

func GetAllServiceNames(w http.ResponseWriter, r *http.Request) {
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

	res, err := q.GetServiceNames(ctx)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	glob, err := json.Marshal(res)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(glob)
}

func PatchServicePassword(w http.ResponseWriter, r *http.Request) {
	var b PatchServicePasswordBody

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
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

	q := services.NewQuerier(conn.Conn())

	_, err = q.PatchServicePassword(ctx, b.Password, b.ServiceName)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteService(w http.ResponseWriter, r *http.Request) {
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

	q := services.NewQuerier(conn.Conn())

	_, err = q.DeleteService(ctx, uuid)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
