package matchmaker

import (
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	predis "github.com/ryank90/matchmaker-sample/pkg/redis"
)

const version string = "alpha-0.0.1"

// Server is the http server instance.
type Server struct {
	srv      *http.Server
	p        *redis.Pool
	sessAddr string
}

// NewServer returns the HTTP server instance.
func NewServer(hostAddr, redisAddr string, sessionAddr string) *Server {
	s := &Server{
		p:        predis.NewPool(redisAddr),
		sessAddr: sessionAddr,
	}

	r := mux.NewRouter()
	r.HandleFunc("/healthz", predis.NewReadinessProbe(s.p))

	s.srv = &http.Server{
		Handler: r,
		Addr:    hostAddr,
	}

	log.Printf("[log][server] connecting to server: %v on port: %v", version, hostAddr)
	log.Printf("[log][server] connecting to redis: %v", redisAddr)
	log.Printf("[log][server] connecting to sessions: %v", sessionAddr)

	return s
}

// Start initialises the server.
func (s *Server) Start() error {
	err := predis.WaitForConnection(s.p)
	if err != nil {
		return errors.Wrap(err, "could not connect to redis")
	}
	return errors.Wrap(s.srv.ListenAndServe(), "error starting the server")
}
