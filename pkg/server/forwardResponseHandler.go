package server

import (
	"net/http"
)

func (l *Listener) handleForwardResponse(code string, w http.ResponseWriter, r *http.Request) {
	responseReference := r.Header.Get("X-RESPONSE-REFERENCE")

	if responseChannel, ok := l.ReferenceToResponse[responseReference]; ok {
		responseChannel <- r

		// Do not responde to the agent until the response is sent to the external requester
		<-l.getResponseSentChannel(responseReference)

	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
