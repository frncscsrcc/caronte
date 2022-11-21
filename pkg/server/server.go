package server

import (
	"log"
	"net/http"
)

type Config struct {
	Host       string
	Port       string
	MaxTimeout int
	BufferSize int
	Secret     string
}

var ServerConfig Config

func Run(config Config) {
	ServerConfig = config
	log.Print("Hook me start")

	listener := NewListener()

	hostAndPort := config.Host + ":" + config.Port
	log.Printf("Listening on %s", hostAndPort)
	if err := http.ListenAndServe(hostAndPort, listener); err != nil {
		panic(err)
	}
}
