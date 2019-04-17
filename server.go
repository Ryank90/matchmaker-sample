package matchmaker

const version string = "alpha-0.0.1"

// Server is the http server instance.
type Server struct {
}

// NewServer returns the HTTP server instance.
func NewServer(hostAddr, redisAddr string, sessionAddr string) *Server {
	return &Server{}
}

// Start initialises the server.
func (s *Server) Start() error {

}
