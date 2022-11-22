package proxy

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (l *Listener) getResponseSentChannel(responseReference string) chan struct{} {
	if responseSentChannel, ok := l.responseSentChannel[responseReference]; ok {
		return responseSentChannel
	}

	newResponseSentChannel := make(chan struct{}, 0)
	l.responseSentChannel[responseReference] = newResponseSentChannel
	return newResponseSentChannel
}

func (l *Listener) handleForwardRequest(code string, w http.ResponseWriter, r *http.Request) {
	log.Print("Received a trigger " + r.URL.Path)

	responseTimeout := getTimeoutChannel(time.Duration(proxyConfig.MaxTimeoutResponse) * time.Second)
	if triggerChannel, exists := l.CodeToPullRequestChannel[code]; !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		responseReference := code + "_" + getRandString(8)
		r.Header.Add("X-RESPONSE-REFERENCE", responseReference)

		responseChannel := make(chan *http.Request, 0)
		l.ReferenceToResponse[responseReference] = responseChannel

		responseSentChannel := l.getResponseSentChannel(responseReference)

		triggerChannel <- r
		log.Printf("Sent request to %s waiting for response [id=%s]", getOriginalPath(r), responseReference)

		select {
		case response := <-responseChannel:
			log.Printf("Received response for %s", responseReference)

			// copy the header
			for name, values := range response.Header {
				for _, value := range values {
					w.Header().Add(name, value)
				}
			}

			// copy the status
			status, err := strconv.Atoi(response.Header.Get("X-STATUS"))
			if err != nil {
				w.WriteHeader(status)
			}

			// copy the body
			io.Copy(w, response.Body)

			responseSentChannel <- struct{}{}

			break
		case <-responseTimeout:
			log.Printf("timeout for %s", responseReference)
			w.WriteHeader(http.StatusRequestTimeout)
			break
		}
	}
}
