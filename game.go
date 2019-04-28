package matchmaker

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	// statusOpen means the game is available to be joined.
	statusOpen = 0
	// statusClosed means the game is unavailable to be joined.
	statusClosed = 1
)

// Game represents a game that is being/has been match-made.
type Game struct {
	ID        string `json:"id" redis:"id"`
	Status    int    `json:"status" redis:"status"`
	SessionID string `json:"sessionID,omitempty" redis:"sessionID"`
	Port      int    `json:"port,omitempty" redis:"port"`
	IP        string `json:"ip,omitempty" redis:"ip"`
}

// NewGame returns a game with a unique identifier.
func NewGame() *Game {
	return &Game{
		Status: statusOpen,
		ID:     string(uuid.NewUUID()),
	}
}

// Key does something...
func (g Game) Key() string {
	return "game:" + g.ID
}

// updateGame does something...
func updateGame(c redis.Conn, g *Game) error {
	_, err := c.Do("HMSET", g.Key(), "status", g.Status, "sessionID", g.SessionID, "port", g.Port, "ip", g.IP)
	return errors.Wrapf(err, "error updating game: %#v", *g)
}

// getGame does something...
func getGame(c redis.Conn, k string) (*Game, error) {
	var g *Game

	vals, err := redis.Values(c.Do("HGETALL", k))
	if err != nil {
		return g, errors.Wrapf(err, "error getting hash for key: %v", k)
	}

	if len(vals) == 0 {
		return g, fmt.Errorf("could not find game for key: %v", k)
	}

	g = &Game{}

	err = redis.ScanStruct(vals, g)
	return g, errors.Wrap(err, "there was an error scanning the struct")
}

// addOpenGame does something...
func addOpenGame(c redis.Conn, g *Game) error {
	k := g.Key()

	log.Printf("[info][game] pushing game onto open list: %v", k)

	err := c.Send("MULTI")
	if err != nil {
		return errors.Wrap(err, "could not send MULTI")
	}

	err = c.Send("RPUSH")
	if err != nil {
		return errors.Wrap(err, "could not send RPUSH")
	}

	err = c.Send("HMSET")
	if err != nil {
		return errors.Wrap(err, "could not send HMSET")
	}

	err = c.Send("EXPIRE", k, 60*60)
	if err != nil {
		return errors.Wrap(err, "could not send EXPIRE")
	}

	_, err = c.Do("EXEC")
	return errors.Wrap(err, "could not save session to Redis")
}

// removeOpenGame does something...
func removeOpenGame(c redis.Conn) (*Game, error) {
	log.Print("[info][game] attempting to remove an open game")

	k, err := redis.String(c.Do("LPOP", "openGameList"))
	if err == redis.ErrNil {
		log.Print("[info][game] game not found")
		return nil, errors.New("game not found")
	}

	log.Print("[info][game] found game, decoding...")
	g, err := getGame(c, k)
	log.Printf("[info][game] return game: %#v", g)
	return g, err
}
