package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/api/middleware"
	v1 "github.com/DillonEnge/service-vault-api/internal/api/v1"
	"github.com/DillonEnge/service-vault-api/internal/bundles"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/sync/errgroup"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Error("Failed to load PORT from env")
		return
	}

	addr := fmt.Sprintf(":%s", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	dburl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		slog.Error("Failed to load DATABASE_URL from env")
		return
	}

	pool, err := pgxpool.Connect(ctx, dburl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /signin", v1.GetSignin)

	mux.HandleFunc("POST /code", v1.PostCode)

	mux.HandleFunc("POST /report", v1.ReportService)
	mux.HandleFunc("DELETE /report", v1.ResetReportedService)

	mux.HandleFunc("GET /admin", v1.GetAdmin)

	mux.HandleFunc("GET /access", v1.GetAccess)
	mux.HandleFunc("POST /access", v1.PostAccess)

	mux.HandleFunc("GET /requests", v1.GetAllRequests)
	mux.HandleFunc("POST /request/code", v1.PostCodeRequest)
	mux.HandleFunc("POST /request/access", v1.PostAccessRequest)
	mux.HandleFunc("DELETE /request", v1.DeleteRequest)

	mux.HandleFunc("GET /services", v1.GetAllServices)
	mux.HandleFunc("GET /service", v1.GetServiceHandler)
	mux.HandleFunc("PATCH /service/password", v1.PatchServicePassword)
	mux.HandleFunc("DELETE /service", v1.DeleteService)

	mux.HandleFunc("GET /policies", v1.GetPolicy)
	mux.HandleFunc("POST /policy", v1.PostPolicy)
	mux.HandleFunc("DELETE /policy", v1.DeletePolicy)

	mux.HandleFunc("GET /groups", v1.GetGroupPolicy)
	mux.HandleFunc("POST /group", v1.PostGroupPolicy)
	mux.HandleFunc("DELETE /group", v1.DeleteGroupPolicy)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]string{"message": "healthy"})
		if err != nil {
			errors.InternalServerError(w, r, err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	b, err := bundles.Build(pool)
	if err != nil {
		slog.Info("Failed to create context bundles")
		return
	}

	s := &http.Server{
		BaseContext:    func(_ net.Listener) context.Context { return ctx },
		Addr:           addr,
		Handler:        middleware.Context(middleware.Logger(middleware.Cors(middleware.Authorizer(mux))), b...),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	g.Go(func() error {
		slog.Info(fmt.Sprintf("Starting server on %s", addr))
		return s.ListenAndServe()
	})

	<-ctx.Done()
	slog.Info("Shutting down...")
	s.Shutdown(ctx)
}
