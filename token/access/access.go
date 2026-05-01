package access

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	persist "github.com/casbin/casbin/v2/persist/string-adapter"
)

// Controller answers authorization questions for a user, system, and action.
//
// Implementations typically evaluate a policy and return whether the permission is
// granted, along with any evaluation/IO error.
type Controller interface {
	// HasAccess reports whether user is allowed to perform action on system.
	// The returned bool indicates whether access is granted. If an error is
	// returned, the boolean result should not be trusted.
	HasAccess(user, system, action string) (bool, error)
}

// NewController constructs a Controller from cfg.
//
// When cfg is nil (disabled), NewController returns (nil, nil).
//
// When enabled, NewController builds a Casbin enforcer using:
//   - the configured model value (cfg.Model) resolved through fs.ReadSource, and
//   - the configured policy value (cfg.Policy) resolved through fs.ReadSource.
func NewController(cfg *Config, fs *os.FS) (Controller, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	config, err := cfg.GetModel(fs)
	if err != nil {
		return nil, err
	}

	policy, err := cfg.GetPolicy(fs)
	if err != nil {
		return nil, err
	}

	model, err := model.NewModelFromString(config)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(model, persist.NewAdapter(policy))
	if err != nil {
		return nil, err
	}

	return &CasbinController{enforcer: enforcer}, nil
}

// CasbinController is a Controller backed by a Casbin enforcer.
//
// It owns the enforcer used to evaluate access checks while exposing only the
// Controller interface used by go-service.
type CasbinController struct {
	enforcer *casbin.Enforcer
}

// HasAccess evaluates whether user is allowed to perform action on system.
//
// It calls the configured Casbin enforcer as:
//
//	c.enforcer.Enforce(user, system, action)
func (c *CasbinController) HasAccess(user, system, action string) (bool, error) {
	return c.enforcer.Enforce(user, system, action)
}
