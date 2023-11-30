package caddyminify

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/tdewolff/minify/v2"
)

type Handler struct {
	Formats []string `json:"formats,omitempty"`

	minifier *minify.M
}

func init() {
	caddy.RegisterModule(&Handler{})
	httpcaddyfile.RegisterHandlerDirective("minify", parseCaddyfile)
}

func (*Handler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.minify",
		New: func() caddy.Module {
			return new(Handler)
		},
	}
}
