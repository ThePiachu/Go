package web

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ThePiachu/Go/Log"
	"google.golang.org/appengine"
)

// Global Mux on which other modules can add route handlers
var Mux = http.NewServeMux()

// Intercepts all requests before giving control to the Mux.
// A panic() in the handler will be caught here where it will
// be logged and a 500 response returned to the client.
func httpInterceptor(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			stackTrace := make([]byte, 1<<16)
			runtime.Stack(stackTrace, false)

			c := appengine.NewContext(r)
			Log.Criticalf(c, "panic caught during request processing: %v, Request: %v\n%s", err, r, stackTrace)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-type", "text/plain; charset=utf-8")
			fmt.Fprint(w, "An error occured Oo..oO\nPlease try again soon.")
		}
	}()

	Mux.ServeHTTP(w, r)
}

func init() {
	http.HandleFunc("/", httpInterceptor)
}
