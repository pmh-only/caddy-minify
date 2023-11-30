package caddyminify

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
	"regexp"
	"slices"
)

func (h *Handler) Provision(_ caddy.Context) error {
	h.minifier = minify.New()

	if h.Formats == nil || slices.Contains(h.Formats, "html") {
		h.minifier.AddFunc("text/html", html.Minify)
	}

	if h.Formats == nil || slices.Contains(h.Formats, "css") {
		h.minifier.AddFunc("text/css", css.Minify)
	}

	if h.Formats == nil || slices.Contains(h.Formats, "svg") {
		h.minifier.AddFunc("image/svg+xml", svg.Minify)
	}

	if h.Formats == nil || slices.Contains(h.Formats, "js") {
		h.minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	}

	if h.Formats == nil || slices.Contains(h.Formats, "json") {
		h.minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	}

	if h.Formats == nil || slices.Contains(h.Formats, "xml") {
		h.minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	}

	print(h.Formats)

	return nil
}

var _ caddy.Provisioner = (*Handler)(nil)
