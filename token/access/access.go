package access

import "github.com/casbin/casbin/v2"

// Controller allows to check different kinds of accesses.
type Controller interface {
	// HasAccess checks if the user can access the system with action.
	HasAccess(user, system, action string) (bool, error)
}

// .NewController for access.
func NewController(cfg *Config) (Controller, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	e, err := casbin.NewEnforcer(cfg.Model, cfg.Policy)
	if err != nil {
		return nil, err
	}

	return &controller{e}, nil
}

type controller struct {
	*casbin.Enforcer
}

func (c *controller) HasAccess(user, system, action string) (bool, error) {
	return c.Enforce(user, system, action)
}
