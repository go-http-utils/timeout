package timeout

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

// Version is this package's version.
var Version = "0.1.0"

// DefaultTimeoutHandler is a convenient timeout handler which
// simply returns "504 Service timeout".
var DefaultTimeoutHandler = http.HandlerFunc(
	func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusGatewayTimeout)
		res.Write([]byte("Service timeout"))
	})

type timeoutWriter struct {
	rw http.ResponseWriter

	status int
	buf    *bytes.Buffer
}

func (tw timeoutWriter) Header() http.Header {
	return tw.rw.Header()
}

func (tw *timeoutWriter) WriteHeader(status int) {
	tw.status = status
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	if tw.status == 0 {
		tw.status = http.StatusOK
	}

	return tw.buf.Write(b)
}

// Handler wraps the http.Handler h with timeout support.
func Handler(h http.Handler, timeout time.Duration, timeoutHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		background := req.Context()
		tCtx, tCancel := context.WithTimeout(background, timeout)
		cCtx, cCancel := context.WithCancel(background)
		req.WithContext(cCtx)

		defer tCancel()

		tw := &timeoutWriter{rw: res, buf: bytes.NewBuffer(nil)}

		go func() {
			h.ServeHTTP(tw, req)
			cCancel()
		}()

		select {
		case <-cCtx.Done():
			res.WriteHeader(tw.status)
			res.Write(tw.buf.Bytes())
		case <-tCtx.Done():
			if err := tCtx.Err(); err == context.DeadlineExceeded {
				cCancel()

				timeoutHandler.ServeHTTP(res, req)
			}
		}
	})
}
