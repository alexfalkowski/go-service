package access

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	persist "github.com/casbin/casbin/v2/persist/string-adapter"
)

// ErrAccessDenied is returned by transport integrations when a policy denies access.
var ErrAccessDenied = errors.New("access: transport policy denied")

// Controller answers authorization questions for request contexts.
//
// Implementations typically evaluate a policy and return whether the permission is
// granted, along with any evaluation/IO error.
type Controller interface {
	// HasAccess reports whether the context's verified user id is allowed to invoke the context's
	// transport service-method.
	// The returned bool indicates whether access is granted. If an error is
	// returned, the boolean result should not be trusted.
	HasAccess(ctx context.Context) (bool, error)
}

// NewController constructs a Controller from cfg.
//
// When cfg is nil (disabled), NewController returns (nil, nil).
//
// When enabled, NewController builds a Casbin enforcer using:
//   - the configured model value (cfg.Model) resolved through [os.FS.ReadSource], and
//   - the configured policy value (cfg.Policy) resolved through [os.FS.ReadSource].
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

// HasAccess evaluates the request context using the controller's configured policy.
//
// The context is expected to contain a verified user id and transport service-method metadata, as populated
// by the standard HTTP and gRPC transport middleware.
func (c *CasbinController) HasAccess(ctx context.Context) (bool, error) {
	return c.enforcer.Enforce(meta.UserID(ctx).Value(), meta.TransportServiceMethod(ctx).Value(), "invoke")
}
