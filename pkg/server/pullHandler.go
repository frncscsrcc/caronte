package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (l *Listener) handlePullRequest(code string, w http.ResponseWriter, r *http.Request) {
	log.Print("Pull request " + r.URL.Path)

	// check if a channel for this client already exists
	if busy, exists := l.Busy[code]; exists && busy {
		fmt.Fprint(w, "busy")
		return
	}

	channel, channelExists := l.CodeToPullRequestChannel[code]
	if !channelExists {
		channel = make(chan *http.Request, ServerConfig.BufferSize)
		l.CodeToPullRequestChannel[code] = channel
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

func waitForData(waitTime time.Duration, triggerChannel <-chan *http.Request, w http.ResponseWriter) {
	timeout := getTimeout(waitTime)

	go func(timeout chan struct{}) {
		time.Sleep(10 * time.Second)
		timeout <- struct{}{}
	}(timeout)

	select {
	case request := <-triggerChannel:
		// copy the headers
		for name, values := range request.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Add custom headers
		w.Header().Add("X-ORIGINAL-PATH", getOriginalPath(request))
		w.Header().Add("X-ORIGINAL-METHOD", request.Method)

		// Copy the body
		io.Copy(w, request.Body)

		break
	case <-timeout:
		w.WriteHeader(http.StatusRequestTimeout)
		break
	}
}
