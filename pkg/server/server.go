package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	DefaultWriteTimeout = 15
	DefaultReadTimeout  = 15
	DefaultIdleTimeout  = 60
)

var srv *http.Server
var r *mux.Router

func init() {
	r = mux.NewRouter()
	r.StrictSlash(true)

	// register handle
	r.HandleFunc("/get", GetHandler).Methods(http.MethodGet)
	r.HandleFunc("/query", GetQueryHandler).Methods(http.MethodGet).Queries("name", "{name}", "age", "{age}")
}

func Serve(addr string, q chan error) {
	srv = &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * DefaultWriteTimeout,
		ReadTimeout:  time.Second * DefaultReadTimeout,
		IdleTimeout:  time.Second * DefaultIdleTimeout,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if q != nil {
				q <- err
			}
		}
	}()
}

func Stop() {
	wait := time.Second * 60
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
}
