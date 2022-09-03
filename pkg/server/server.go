package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

type ForwardedRequest struct {
	Method       string
	Body         []byte
	Headers      []http.Header
	OriginalPath string
}

type Listener struct {
	CodeToClient map[string]chan ForwardedRequest
	Busy         map[string]bool
}

func NewListener() *Listener {
	return &Listener{
		CodeToClient: make(map[string]chan ForwardedRequest),
		Busy:         make(map[string]bool),
	}
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "wrong request")
		return
	}

	action := parts[1]
	code := parts[2]
	if action == "receive" {
		l.receive(code, w, r)
		return
	} else if action == "trigger" {
		l.trigger(code, w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "wrong request")

}

func (l *Listener) receive(code string, w http.ResponseWriter, r *http.Request) {
	log.Print("Receive " + r.URL.Path)
	// check if a channel for this client already exists
	if busy, exists := l.Busy[code]; exists && busy {
		fmt.Fprint(w, "busy")
	} else {
		channel, channelExists := l.CodeToClient[code]
		if !channelExists {
			channel = make(chan ForwardedRequest, ServerConfig.BufferSize)
			l.CodeToClient[code] = channel
		}
		l.Busy[code] = true
		defer func() { l.Busy[code] = false }()

		waitTime := time.Duration(ServerConfig.MaxTimeout) * time.Second
		if headerTimeoutString := r.Header.Get("X-TIMEOUT"); headerTimeoutString != "" {
			if timeoutInt, err := strconv.Atoi(headerTimeoutString); err == nil {
				waitTime = time.Duration(timeoutInt) * time.Second
			}
		}
		waitForData(waitTime, channel, w)
	}
}

func getTimeout(d time.Duration) chan struct{} {
	timeout := make(chan struct{})
	go func() {
		time.Sleep(d)
		timeout <- struct{}{}
	}()
	return timeout
}

func waitForData(waitTime time.Duration, channel <-chan ForwardedRequest, w http.ResponseWriter) {
	timeout := getTimeout(waitTime)

	go func(timeout chan struct{}) {
		time.Sleep(10 * time.Second)
		timeout <- struct{}{}
	}(timeout)

	select {
	case request := <-channel:
		for _, header := range request.Headers {
			for name, values := range header {
				for _, value := range values {
					w.Header().Add(name, value)
				}
			}
		}
		w.Header().Add("X-ORIGINAL-PATH", request.OriginalPath)
		w.Header().Add("X-ORIGINAL-METHOD", request.Method)
		io.Copy(w, bytes.NewBuffer(request.Body))

		break
	case <-timeout:
		w.WriteHeader(http.StatusRequestTimeout)
		break
	}
}

func (l *Listener) trigger(code string, w http.ResponseWriter, r *http.Request) {
	log.Print("Trigger " + r.URL.Path)

	timeout := getTimeout(2 * time.Second)
	if channel, exists := l.CodeToClient[code]; !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
		r.Body.Close()

		originalPath := "/" + strings.Join(
			strings.Split(r.URL.Path, "/")[3:],
			"/",
		)

		headers := make([]http.Header, 0)
		for name, values := range r.Header {
			headers = append(headers, http.Header{name: values})
		}

		request := ForwardedRequest{
			Headers:      headers,
			OriginalPath: originalPath,
			Method:       r.Method,
			Body:         body,
		}

		select {
		case channel <- request:
			w.WriteHeader(http.StatusAccepted)
			break
		case <-timeout:
			w.WriteHeader(http.StatusRequestTimeout)
			break
		}

	}
}
