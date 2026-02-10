package access

// ModelConfig is the embedded Casbin RBAC model used by this package.
//
// It models requests as (sub, obj, act) and uses role inheritance via `g`.
// Policies are evaluated with an "allow" effect and matched by exact object and action.
const ModelConfig = `
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
