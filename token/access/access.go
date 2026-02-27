package access

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	file "github.com/casbin/casbin/v2/persist/file-adapter"
)

// Controller answers authorization questions of the form “is user allowed to do X?”
//
// Implementations typically evaluate a policy and return whether the permission is
// granted, along with any evaluation/IO error.
type Controller interface {
	// HasAccess reports whether user is allowed the given permission.
	//
	// permission is expected to be in the form:
	//
	//	<system>:<action>
	//
	// For example:
	//
	//	"service:read"
	//
	// Implementations split this string into object (system) and action components
	// and evaluate a policy. The returned bool indicates whether access is granted.
	// If an error is returned, the boolean result should not be trusted.
	HasAccess(user, permission string) (bool, error)
}

// NewController constructs a Controller from cfg.
//
// When cfg is nil (disabled), NewController returns (nil, nil).
//
// When enabled, NewController builds a Casbin enforcer using:
//   - the embedded RBAC model definition (ModelConfig), and
//   - the configured policy value (cfg.Policy) via Casbin’s file adapter.
//
// Note: the policy string is passed directly to the adapter constructor. Ensure
// cfg.Policy matches what the underlying adapter expects in your environment
// (for example a path vs. a literal policy payload).
//
// Any model parse error triggers a panic via runtime.Must; this is treated as a
// programmer error because ModelConfig is a package constant.
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

// CasbinController is a Controller backed by a Casbin enforcer.
//
// It embeds *casbin.Enforcer so callers can access Casbin capabilities directly
// when needed, while still satisfying the Controller interface used by go-service.
type CasbinController struct {
	*casbin.Enforcer
}

// HasAccess evaluates whether user is allowed the given permission.
//
// It splits permission on the first ":" into (system, action) and calls the embedded
// Casbin enforcer as:
//
//	c.Enforce(user, system, action)
//
// The permission string is expected to be in the form "<system>:<action>". If the
// string does not contain ":", the action will be empty and the policy will be
// evaluated accordingly.
func (c *CasbinController) HasAccess(user, permission string) (bool, error) {
	system, action := strings.CutColon(permission)
	return c.Enforce(user, system, action)
}
