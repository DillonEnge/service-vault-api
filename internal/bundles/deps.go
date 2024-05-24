package bundles

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/DillonEnge/service-vault-api/internal/api/middleware"
	"github.com/DillonEnge/service-vault-api/internal/authz"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/casbin/casbin/v2"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/cychiuae/casbin-pg-adapter"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrEnvarNotFound = errors.New("Envar not found")
)

func Build(pool *pgxpool.Pool) ([]*middleware.ContextBundle, error) {
	endpoint, ok := os.LookupEnv("CASDOOR_ENDPOINT")
	if !ok {
		slog.Error("Failed to load CASDOOR_ENDPOINT from env")
		return nil, ErrEnvarNotFound
	}

	clientId, ok := os.LookupEnv("CASDOOR_CLIENT_ID")
	if !ok {
		slog.Error("Failed to load CASDOOR_CLIENT_ID from env")
		return nil, ErrEnvarNotFound
	}

	clientSecret, ok := os.LookupEnv("CASDOOR_CLIENT_SECRET")
	if !ok {
		slog.Error("Failed to load CASDOOR_CLIENT_SECRET from env")
		return nil, ErrEnvarNotFound
	}

	certFileBytes, err := os.ReadFile("./.cert")
	if err != nil {
		slog.Error("Failed to read cert file")
		return nil, err
	}

	cert := string(certFileBytes)

	organizationName, ok := os.LookupEnv("CASDOOR_ORGANIZATION_NAME")
	if !ok {
		slog.Error("Failed to load CASDOOR_ORGANIZATION_NAME from env")
		return nil, ErrEnvarNotFound
	}

	applicationName, ok := os.LookupEnv("CASDOOR_APPLICATION_NAME")
	if !ok {
		slog.Error("Failed to load CASDOOR_APPLICATION_NAME from env")
		return nil, ErrEnvarNotFound
	}

	authConfig := &casdoorsdk.AuthConfig{
		Endpoint:         endpoint,
		ClientId:         clientId,
		ClientSecret:     clientSecret,
		Certificate:      cert,
		OrganizationName: organizationName,
		ApplicationName:  applicationName,
	}

	casdoorClient := casdoorsdk.NewClientWithConf(authConfig)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	tableName := "casbin"
	adapter, err := casbinpgadapter.NewAdapter(db, tableName)
	// If you are using db schema
	// myDBSchema := "mySchema"
	// adapter, err := casbinpgadapter.NewAdapterWithDBSchema(db, myDBSchema, tableName)
	if err != nil {
		slog.Error("Failed to create adapter")
		return nil, err
	}

	// policyPath := "dev/policy.csv"
	modelPath := "dev/model.conf"

	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		slog.Error("Failed to create new enforcer")
		return nil, err
	}

	e.AddGroupingPolicy("dillon.enge@gmail.com", "admin")

	rules := [][]string{
		{"admin", "/groups", "GET"},
		{"admin", "/group", "POST"},
		{"admin", "/group", "DELETE"},
		{"admin", "/policy", "DELETE"},
		{"admin", "/policy", "POST"},
		{"admin", "/policies", "GET"},
	}

	for _, rule := range rules {
		_, err = e.AddPolicy(rule)
		if err != nil {
			slog.Error("Failed to add policies")
			return nil, err
		}
	}

	casbinClient := &authz.Client{
		Enforcer: e,
	}

	b := []*middleware.ContextBundle{
		{
			Key: consts.CtxKeyDbPool,
			Val: pool,
		},
		{
			Key: consts.CtxKeyCasdoorClient,
			Val: casdoorClient,
		},
		{
			Key: consts.CtxKeyCasbinClient,
			Val: casbinClient,
		},
	}

	return b, nil
}
