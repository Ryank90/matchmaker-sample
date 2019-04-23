package matchmaker

import "github.com/garyburd/redigo/redis"

const (
	maxRetries = 20
)

// Session represents a game session.
type Session struct {
	ID   string `json:""`
	Port int    `json:""`
	IP   string `json:""`
}

func (s *Server) generateSessionForGame(c redis.Conn, g *Game) (*Game, error) {

}

func (s *Server) sessionIPAndPort(sess Session) (Session, error) {

}
