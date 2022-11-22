package proxy

import (
	"log"
	"net/http"
)

type Config struct {
	Host               string
	Port               string
	MaxTimeoutLongPoll int
	MaxTimeoutResponse int
	BufferSize         int
	Secret             string
}

var proxyConfig Config

func Run(config Config) {
	proxyConfig = config
	log.Print("Hook me start")

	listener := NewListener()

	hostAndPort := config.Host + ":" + config.Port
	log.Printf("Proxy listening on %s", hostAndPort)
	if err := http.ListenAndServe(hostAndPort, listener); err != nil {
		panic(err)
	}
}
