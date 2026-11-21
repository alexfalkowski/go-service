package logger

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v3/meta"
	"github.com/alexfalkowski/go-service/v3/net/http"
	httpmeta "github.com/alexfalkowski/go-service/v3/net/http/meta"
	"github.com/alexfalkowski/go-service/v3/net/url"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-service/v3/telemetry/logger"
)

func serverRequestLog(req *http.Request, policy *http.RoutePolicy) (string, []logger.Attr) {
	path := cmp.Or(req.URL.EscapedPath(), "/")
	attrs := requestAttrs(path)
	route := routePath(req.Pattern)
	if route != "" {
		attrs = append(attrs, logger.String("http.route", route))
	}
	target := cmp.Or(route, path)

	var rpcPath string
	if policy.IsRPC(req) {
		rpcPath = target
	}

	return requestLog(req.Method, target, rpcPath, attrs)
}

func clientRequestLog(req *http.Request) (string, []logger.Attr) {
	path := cmp.Or(req.URL.EscapedPath(), "/")
	attrs := requestAttrs(path)
	if host := req.URL.Hostname(); host != "" {
		attrs = append(attrs, logger.String("server.address", host))
	}

	// Redirects inherit the original request context and carry the response that redirected them.
	var rpcPath string
	if req.Response == nil {
		if template := httpmeta.URLTemplate(req.Context()); template != "" {
			attrs = append(attrs, logger.String("url.template", template))
			path = template
		}
		rpcPath = httpmeta.RPCMethod(req.Context())
	}

	return requestLog(req.Method, strings.Concat(req.URL.Host, path), rpcPath, attrs)
}

func requestAttrs(path string) []logger.Attr {
	attrs := make([]logger.Attr, 0, 9)
	return append(attrs, logger.String(meta.SystemKey, "http"), logger.String("url.path", path))
}

func requestLog(verb, target, rpcPath string, attrs []logger.Attr) (string, []logger.Attr) {
	verb = strings.ToLower(verb)
	attrs = append(attrs, logger.String("http.request.method", verb))

	if service, method, ok := url.SplitPath(rpcPath); ok {
		attrs = append(attrs, logger.String(meta.ServiceKey, service), logger.String(meta.MethodKey, method))
		return strings.Concat("http: ", service, "/", method), attrs
	}

	return strings.Concat("http: ", verb, " ", target), attrs
}

func routePath(pattern string) string {
	prefix, _, ok := strings.Cut(pattern, "/")
	if !ok {
		return ""
	}

	return pattern[len(prefix):]
}
