package timeout_test

import (
	"net/http"
	"time"

	"github.com/go-http-utils/timeout"
)

func Example() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", timeout.Handler(mux, time.Second*10, timeout.DefaultTimeoutHandler))
}
