package caddyminify

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
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

	if recorder.Header().Get("Content-Encoding") != "" {
		w.WriteHeader(recorder.Status())
		_, err = w.Write(resBuffer.Bytes())
		return err
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

var _ caddyhttp.MiddlewareHandler = (*Handler)(nil)
