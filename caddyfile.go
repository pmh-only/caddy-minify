package caddyminify

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func (h *Handler) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextArg() {
			h.Formats = append(h.Formats, d.Val())
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			if d.Val() == "formats" {
				for d.NextArg() {
					h.Formats = append(h.Formats, d.Val())
				}
			}
		}
	}

	return nil
}

func parseCaddyfile(helper httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var h Handler
	err := h.UnmarshalCaddyfile(helper.Dispenser)
	return &h, err
}
