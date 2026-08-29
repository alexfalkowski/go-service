# 🎫 Tokens

[← Back to README](../README.md)

Token configuration is rooted at `token.Config`, usually nested under transport config as `transport.http.token` and/or `transport.grpc.token` (via the shared server-side transport config).

Supported token `kind` values:

- `jwt`
- `paseto`

## Access control (Casbin)

Access control is configured once at the transport level and shared by all
enabled HTTP and gRPC server stacks:

```yaml
transport:
  access:
    model: file:./config/rbac.conf
    policy: file:./config/rbac.csv
```

When `access` is configured, the standard HTTP and gRPC server stacks enforce
the policy after token authentication and before application handlers run. Omit
`access` to leave transport authorization disabled.

The model is based on Casbin RBAC:
<https://github.com/casbin/casbin/blob/master/examples/rbac_model.conf>

Example `rbac.conf`:

```ini
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
```

Policies use the verified user id as `sub`, `meta.TransportServiceMethod` as
`obj`, and `invoke` as `act`. Example `rbac.csv`:

```csv
p, reader, http:GET /users/{id}, invoke
p, writer, http:POST /users, invoke
p, greeter, grpc:/greet.v1.GreeterService/SayHello, invoke
g, frontend, reader
g, admin, reader
g, admin, writer
g, billing-service, greeter
```

The `p` rows define permissions and must match the model's `p = sub, obj, act`
shape, so they include `invoke`. The `g` rows define role membership and match
`g = _, _`, so they only contain `subject, role`.

> [!WARNING]
> Casbin's string policy adapter can skip malformed policy rows without failing
> startup. Validate policy files before deployment; a successful startup does
> not prove that every configured row was loaded.

For HTTP servers the object uses the matched route pattern when available, such
as `http:GET /users/{id}`. HTTP tokens are authenticated against the concrete
request method and path, such as `GET /users/123`; access policy enforcement
uses the canonical route pattern. gRPC tokens are authenticated against the full
method name, such as `/greet.v1.GreeterService/SayHello`; access policy
enforcement uses the transport service-method object, such as
`grpc:/greet.v1.GreeterService/SayHello`.

> [!NOTE]
> `access.model` and `access.policy` are resolved through `os.FS.ReadSource`; use `file:` for files, `env:` for environment-provided content, or literal content.
>
> Access config builds an injectable controller for authorization checks. The built-in HTTP and gRPC server stacks authenticate tokens, store the verified user id, and enforce the configured Casbin policy before application handlers run.

## JWT

JWT config:

```yaml
transport:
  http:
    token:
      kind: jwt
      jwt:
        iss: my-service
        exp: 1h
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/ed25519.pub
            private: file:/keys/ed25519
          old:
            public: file:/keys/ed25519-old.pub
```

Important behavior:

- JWT generation signs with `jwt.key`; verification requires the token `kid` header to select an entry in `jwt.keys`.
- `exp` is parsed as a Go duration string; invalid values can fail fast.
- `leeway` is optional clock-skew tolerance for verification; keep it small because it extends acceptance around `iat`/`nbf` and `exp`.

> [!IMPORTANT]
> JWT generation and verification use Ed25519 key material from `jwt.keys`. Keep private key material only on services that mint tokens; verifiers only need public keys.

All token `exp` and non-zero `leeway` values are Go duration strings and must be positive whole-second durations. Values such
as `1s`, `15m`, and `24h` validate; sub-second values such as `500ms` do not.

## Paseto

Paseto config:

```yaml
transport:
  http:
    token:
      kind: paseto
      paseto:
        iss: my-service
        exp: 1h
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/ed25519.pub
            private: file:/keys/ed25519
          old:
            public: file:/keys/ed25519-old.pub
```

> [!NOTE]
> The PASETO implementation issues **v4 public** tokens. Generation signs with `paseto.key`, writes that id as footer `kid`, and verification selects the public key from `paseto.keys`. `paseto.leeway` is optional clock-skew tolerance for verification.
