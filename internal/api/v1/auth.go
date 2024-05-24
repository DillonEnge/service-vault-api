package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DillonEnge/service-vault-api/internal/api/errors"
	"github.com/DillonEnge/service-vault-api/internal/authz"
	"github.com/DillonEnge/service-vault-api/internal/consts"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type Policy struct {
	Sub string `json:"sub"`
	Obj string `json:"obj"`
	Act string `json:"act"`
}

type Group struct {
	Group string `json:"group"`
	Sub   string `json:"sub"`
}

func GetSignin(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	client := r.Context().Value(consts.CtxKeyCasdoorClient).(*casdoorsdk.Client)

	token, err := client.GetOAuthToken(code, state)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	glob, err := json.Marshal(map[string]interface{}{
		"access_token": token.AccessToken,
		"token_type":   token.TokenType,
	})
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(glob)
}

func GetPolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cbinClient := ctx.Value(consts.CtxKeyCasbinClient).(*authz.Client)

	policy := cbinClient.Enforcer.GetPolicy()

	b, err := json.Marshal(policy)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func PostPolicy(w http.ResponseWriter, r *http.Request) {
	var p Policy

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	cbinClient := ctx.Value(consts.CtxKeyCasbinClient).(*authz.Client)

	ok, err := cbinClient.Enforcer.AddPolicy(p.Sub, p.Obj, p.Act)
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

	w.WriteHeader(http.StatusOK)
}

func DeletePolicy(w http.ResponseWriter, r *http.Request) {
	sub := r.URL.Query().Get("sub")
	if sub == "" {
		slog.Error("Failed to find mandatory 'sub' param in query string")
		errors.InternalServerError(w, r, nil)
		return
	}

	obj := r.URL.Query().Get("obj")
	if obj == "" {
		slog.Error("Failed to find mandatory 'obj' param in query string")
		errors.InternalServerError(w, r, nil)
		return
	}

	act := r.URL.Query().Get("act")
	if act == "" {
		slog.Error("Failed to find mandatory 'act' param in query string")
		errors.InternalServerError(w, r, nil)
		return
	}

	ctx := r.Context()

	cbinClient := ctx.Value(consts.CtxKeyCasbinClient).(*authz.Client)

	ok, err := cbinClient.Enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	if !ok {
		slog.Error("Policy rule not present")
		errors.InternalServerError(w, r, nil)
		return
	}

	err = cbinClient.Enforcer.SavePolicy()
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetGroupPolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cbinClient := ctx.Value(consts.CtxKeyCasbinClient).(*authz.Client)

	policy := cbinClient.Enforcer.GetGroupingPolicy()

	b, err := json.Marshal(policy)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func PostGroupPolicy(w http.ResponseWriter, r *http.Request) {
	var g Group

	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	cbinClient := ctx.Value(consts.CtxKeyCasbinClient).(*authz.Client)

	ok, err := cbinClient.Enforcer.AddGroupingPolicy(g.Group, g.Sub)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	if !ok {
		slog.Error("Grouping policy rule already present")
		errors.InternalServerError(w, r, nil)
		return
	}

	err = cbinClient.Enforcer.SavePolicy()
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteGroupPolicy(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	if group == "" {
		slog.Error("Failed to find mandatory 'group' param in query string")
		errors.InternalServerError(w, r, nil)
		return
	}

	sub := r.URL.Query().Get("sub")
	if sub == "" {
		slog.Error("Failed to find mandatory 'sub' param in query string")
		errors.InternalServerError(w, r, nil)
		return
	}

	ctx := r.Context()

	cbinClient := ctx.Value(consts.CtxKeyCasbinClient).(*authz.Client)

	ok, err := cbinClient.Enforcer.RemoveGroupingPolicy(group, sub)
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}
	if !ok {
		slog.Error("Grouping policy rule not present")
		errors.InternalServerError(w, r, nil)
		return
	}

	err = cbinClient.Enforcer.SavePolicy()
	if err != nil {
		errors.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
