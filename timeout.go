package timeout

import (
	"context"
	"net/http"
	"time"
)

// Version is this package's version.
var Version = "0.0.1"

// DefaultTimeoutHandler is a convenient timeout handler whose behaviour is
// simply return "504 Service Timeout".
var DefaultTimeoutHandler = http.HandlerFunc(
	func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusGatewayTimeout)
		res.Write([]byte("Service Timeout"))
	})

// Handler wraps the http.Handler h with timeout support.
func Handler(h http.Handler, timeout time.Duration, timeoutHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		background := req.Context()
		tCtx, tCancel := context.WithTimeout(background, timeout)
		cCtx, cCancel := context.WithCancel(background)
		req.WithContext(cCtx)

		defer tCancel()
		defer cCancel()

		go h.ServeHTTP(res, req)

		select {
		case <-cCtx.Done():
			return
		case <-tCtx.Done():
			if err := tCtx.Err(); err == context.DeadlineExceeded {
				cCancel()

				timeoutHandler.ServeHTTP(res, req)
			}
		}
	})
}
