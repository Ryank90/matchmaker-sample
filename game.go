package matchmaker

import "k8s.io/apimachinery/pkg/util/uuid"

const (
	// statusOpen means the game is available to be joined.
	statusOpen = 0
	// statusClosed means the game is unavailable to be joined.
	statusClosed = 1
)

// Game represents a game that is being/has been match-made.
type Game struct {
	ID        string `json:"" redis:""`
	Status    int    `json:"" redis:""`
	SessionID string `json:"" redis:""`
	Port      int    `json:"" redis:""`
	IP        string `json:"" redis:""`
}

// NewGame returns a game with a unique identifier.
func NewGame() *Game {
	return &Game{
		Status: statusOpen,
		ID:     string(uuid.NewUUID()),
	}
}
