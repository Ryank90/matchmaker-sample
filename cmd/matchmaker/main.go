package main

import (
	"log"
	"os"

	"github.com/Ryank90/matchmaker-sample"
)

func main() {
	log.Print("[info][main] creating server...")
	s := matchmaker.NewServer(":"+os.Getenv("PORT"), os.Getenv("REDIS_SERVICE"), os.Getenv("SESSION_SERVICE"))
	if err := s.Start(); err != nil {
		log.Fatalf("[error][main] %+v", err)
	}
}
