package caddyminify

import (
	"bytes"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

func init() {
	caddy.RegisterModule(&Handler{})
	httpcaddyfile.RegisterHandlerDirective("minify", parseCaddyfile)
}

type Handler struct {
	minifier *minify.M
}

func (*Handler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.minify",
		New: func() caddy.Module {
			return new(Handler)
		},
	}
}

func (h *Handler) Provision(_ caddy.Context) error {
	h.minifier = minify.New()
	h.minifier.AddFunc("text/html", html.Minify)
	h.minifier.AddFunc("text/css", css.Minify)
	h.minifier.AddFunc("image/svg+xml", svg.Minify)
	h.minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	h.minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	h.minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	resBuffer := bufferPool.Get().(*bytes.Buffer)
	resBuffer.Reset()
	defer bufferPool.Put(resBuffer)

	shouldBuffer := func(_ int, _ http.Header) bool { return true }
	recorder := caddyhttp.NewResponseRecorder(w, resBuffer, shouldBuffer)

	err := next.ServeHTTP(recorder, r)
	if err != nil {
		return err
	}

	if resBuffer.Len() < 1 {
		w.WriteHeader(recorder.Status())
		return nil
	}

	result := &bytes.Buffer{}
	contentType := recorder.Header().Get("Content-Type")

	err = h.minifier.Minify(contentType, result, resBuffer)
	if err != nil {
		w.WriteHeader(recorder.Status())
		_, err = w.Write(resBuffer.Bytes())
		return err
	}

	w.Header().Set("Content-Length", strconv.Itoa(result.Len()))
	w.WriteHeader(recorder.Status())
	_, err = w.Write(result.Bytes())

	return err
}

func parseCaddyfile(_ httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Handler
	return &m, nil
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var (
	_ caddy.Provisioner           = (*Handler)(nil)
	_ caddyhttp.MiddlewareHandler = (*Handler)(nil)
)
