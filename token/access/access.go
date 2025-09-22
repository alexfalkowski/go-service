package access

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	file "github.com/casbin/casbin/v2/persist/file-adapter"
)

// Controller allows to check different kinds of accesses.
type Controller interface {
	// HasAccess checks if the user can access the system with action.
	HasAccess(user, permission string) (bool, error)
}

// NewController for access.
func NewController(cfg *Config) (Controller, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	m, err := model.NewModelFromString(ModelConfig)
	runtime.Must(err)

	e, err := casbin.NewEnforcer(m, file.NewAdapter(cfg.Policy))
	if err != nil {
		return nil, err
	}

	return &CasbinController{e}, nil
}

// CasbinController for access.
type CasbinController struct {
	*casbin.Enforcer
}

// HasAccess just calls Enforce.
func (c *CasbinController) HasAccess(user, permission string) (bool, error) {
	system, action := strings.CutColon(permission)
	return c.Enforce(user, system, action)
}
