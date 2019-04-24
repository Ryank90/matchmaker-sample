package matchmaker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

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
	p := s.sessAddr + "/session"

	r, err := http.Post(p, "application/json", nil)
	if err != nil {
		return g, errors.Wrap(err, "error calling session")
	}

	defer r.Body.Close()

	sess := Session{}
	err = json.NewDecoder(r.Body).Decode(&sess)
	if err != nil {
		return g, errors.WithStack(err)
	}

	g.SessionID = sess.ID
	sess, err = s.sessionIPAndPort(sess)
	g.Port = sess.Port
	g.IP = sess.IP
	g.Status = statusClosed
}

func (s *Server) sessionIPAndPort(sess Session) (Session, error) {
	var body io.ReadCloser

	for i := 0; i <= maxRetries; i++ {
		r := s.sessAddr + "/session/" + url.QueryEscape(sess.ID)
		res, err := http.Get(r)
		if err != nil {
			return sess, errors.Wrap(err, "error getting session information")
		}

		if r.StatusCode == http.StatusOK {
			log.Printf("[info][session] recieved session data, status: %v", r.StatusCode)
			body = r.Body
			break
		}

		err = r.Body.Close()
		if err != nil {
			log.Printf("[warn][session] could not close body: %v", err)
		}

		log.Printf("[info][session] session: %v data could not be found, please try again", sess.ID)

		time.Sleep(time.Second)
	}

	defer body.Close()

	if body == nil {
		return sess, errors.Errorf("could not get session: %v", sess.ID)
	}

	return sess, errors.Wrap(json.NewDecoder(body).Decode(&sess), "could not decode json to session")
}
