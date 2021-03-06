package redis

import (
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"

	"github.com/garyburd/redigo/redis"
)

// WaitForConnection pings Redis until a connection can be made.
func WaitForConnection(p *redis.Pool) error {
	return backoff.Retry(func() error {
		c := p.Get()
		defer c.Close()

		_, err := c.Do("PING")
		if err != nil {
			log.Printf("[warn][redis] retrying connection to Redis: %v", err)
		} else {
			log.Print("[info][redis] connected successfully.")
		}

		return errors.Wrap(err, "could not connect to Redis")

	}, backoff.NewExponentialBackOff())
}

// NewPool returns a new Redis pool.
func NewPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 4 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return errors.Wrap(err, "error with ping on testonborrow")
		},
	}
}

// NewReadinessProbe creates a HTTP readiness probe.
func NewReadinessProbe(p *redis.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := p.Get()
		defer conn.Close()

		_, err := conn.Do("PING")
		if err != nil {
			log.Printf("[warn][redis] readiness probe failed. %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
