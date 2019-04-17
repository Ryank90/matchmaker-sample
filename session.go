package matchmaker

const (
	maxRetries = 20
)

// Session represents a game session.
type Session struct {
	ID   string `json:""`
	Port int    `json:""`
	IP   string `json:""`
}
