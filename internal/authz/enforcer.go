package authz

import (
	"github.com/casbin/casbin/v2"
)

type Client struct {
	Enforcer *casbin.Enforcer
}

func (c *Client) GetAuthorization(sub string, obj string, act string) bool {
	if res, _ := c.Enforcer.Enforce(sub, obj, act); res {
		return true
	}

	return false
}
