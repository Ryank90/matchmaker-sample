package matchmaker

import (
	"encoding/json"
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

type handler func(*Server, http.ResponseWriter, *http.Request) error

// NewServer returns the HTTP server instance.
func NewServer(hostAddr, redisAddr string, sessionAddr string) *Server {
	s := &Server{
		p:        predis.NewPool(redisAddr),
		sessAddr: sessionAddr,
	}

	r := s.routes()

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

func (s *Server) routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", predis.NewReadinessProbe(s.p))
	r.HandleFunc("/game", s.middleware(gameHandler)).Methods("POST")
	r.HandleFunc("/game/{id}", s.middleware(getHandler)).Methods("GET")

	return r
}

func (s *Server) middleware(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(s, w, r)
		if err != nil {
			log.Printf("[error][server] %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func gameHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	c := s.p.Get()
	defer c.Close()

	log.Print("[info][route] match to a game")
	g, err := removeOpenGame(c)

	if err != nil {
		if err != errors.New("game not found") {
			return err
		}

		g = NewGame()
		err := addOpenGame(c, g)

		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		g, err = s.generateSessionForGame(c, g)
		if err != nil {
			return err
		}

		err := updateGame(c, g)
		if err != nil {
			return err
		}
	}

	return errors.Wrap(json.NewEncoder(w).Encode(g), "error encoding to json")
}

func getHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	c := s.p.Get()
	defer c.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("[info][route] retrieving game: %v", id)
	g, err := getGame(c, "")
	if err != nil {
		return err
	}

	return errors.Wrap(json.NewEncoder(w).Encode(g), "error encoding game to json")
}
