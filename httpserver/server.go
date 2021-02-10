package httpserver

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/valyala/fasthttp/reuseport"
)

type server struct {
	HTTPServer *fasthttp.Server
	Router     *router.Router
	ListenAddr string
	ServerName string
}

// NewServer creates a new HTTP Serverg
func NewServer(listenaddr string, servername string) *server {

	// define router
	r := router.New()

	return &server{
		HTTPServer: newHTTPServer(errorCatcher(r.Handler), servername),
		Router:     r,
		ListenAddr: listenaddr,
		ServerName: servername,
	}
}

// NewServer creates a new HTTP Server
// TODO: configuration should be configurable
func newHTTPServer(h fasthttp.RequestHandler, servername string) *fasthttp.Server {
	return &fasthttp.Server{
		Handler:              h,
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         10 * time.Second,
		MaxConnsPerIP:        500,
		MaxRequestsPerConn:   500,
		MaxKeepaliveDuration: 5 * time.Second,
		Name:                 servername,
		ReduceMemoryUsage:    true,
	}
}

// Run starts the HTTP server and performs a graceful shutdown
func (s *server) Run() {
	// NOTE: Package reuseport provides a TCP net.Listener with SO_REUSEPORT support.
	// SO_REUSEPORT allows linear scaling server performance on multi-CPU servers.

	// create a fast listener ;)
	ln, err := reuseport.Listen("tcp4", s.ListenAddr)
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	// create a graceful shutdown listener
	duration := 5 * time.Second
	graceful := NewGracefulListener(ln, duration)

	// Error handling
	listenErr := make(chan error, 1)

	/// Run server
	go func() {
		log.Printf("%s - Server starting on port %v", s.ServerName, graceful.Addr())
		//log.Printf("%s - Press Ctrl+C to stop", s.ServerName)
		// listenErr <- s.HTTPServer.ListenAndServe(":" + cfg.Port)
		listenErr <- s.HTTPServer.Serve(graceful)

	}()

	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Handle channels/graceful shutdown
	for {
		select {
		// If server.ListenAndServe() cannot start due to errors such
		// as "port in use" it will return an error.
		case err := <-listenErr:
			if err != nil {
				log.Fatalf("listener error: %s", err)
			}
			os.Exit(0)
		// handle termination signal
		case <-osSignals:
			fmt.Printf("\n")
			log.Printf("%s - Shutdown signal received.\n", s.ServerName)

			// Servers in the process of shutting down should disable KeepAlives
			// FIXME: This causes a data race
			s.HTTPServer.DisableKeepalive = true

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			if err := graceful.Close(); err != nil {
				log.Fatalf("error with graceful close: %s", err)
			}

			log.Printf("%s - Server gracefully stopped.\n", s.ServerName)
		}
	}
}
