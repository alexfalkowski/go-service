package access

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	file "github.com/casbin/casbin/v2/persist/file-adapter"
)

// Controller checks whether a subject has permission to perform an action.
type Controller interface {
	// HasAccess reports whether user is allowed the given permission.
	//
	// permission is expected to be in the form "system:action".
	HasAccess(user, permission string) (bool, error)
}

// NewController constructs a Controller from cfg.
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

// CasbinController is a Controller backed by Casbin.
type CasbinController struct {
	*casbin.Enforcer
}

// HasAccess enforces permission for user using the embedded Casbin enforcer.
func (c *CasbinController) HasAccess(user, permission string) (bool, error) {
	system, action := strings.CutColon(permission)
	return c.Enforce(user, system, action)
}
