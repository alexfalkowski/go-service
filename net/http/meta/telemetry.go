package meta

import "github.com/alexfalkowski/go-service/v3/context"

const (
	rpcMethodKey   = context.Key("http-rpc-method")
	urlTemplateKey = context.Key("http-url-template")
)

// WithRPCMethod attaches a full RPC method such as "/greet.v1.Greeter/SayHello" to the context of
// an outbound HTTP request. Pass the derived context to that request to name it in the client logger.
func WithRPCMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, rpcMethodKey, method)
}

// RPCMethod returns the full outbound HTTP RPC method in ctx, or an empty string if unset.
func RPCMethod(ctx context.Context) string {
	method, _ := ctx.Value(rpcMethodKey).(string)
	return method
}

// WithURLTemplate attaches a path template such as "/repos/{owner}/{repo}" to the context of one
// outbound HTTP request. The client logger uses it to name the initial request and its retries.
// Use a stable template with placeholders for variable segments.
func WithURLTemplate(ctx context.Context, template string) context.Context {
	return context.WithValue(ctx, urlTemplateKey, template)
}

// URLTemplate returns the outbound HTTP path template in ctx, or an empty string if unset.
func URLTemplate(ctx context.Context) string {
	template, _ := ctx.Value(urlTemplateKey).(string)
	return template
}
