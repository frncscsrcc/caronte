package proxy

import (
	"fmt"
	"net/http"
	"strings"
)

type Listener struct {
	CodeToPullRequestChannel map[string]chan *http.Request
	Busy                     map[string]bool
	ReferenceToResponse      map[string]chan *http.Request
	responseSentChannel      map[string]chan struct{}
}

func NewListener() *Listener {
	return &Listener{
		// Map the resource code with the request received by the external requester
		CodeToPullRequestChannel: make(map[string]chan *http.Request),

		// Mark a resource as busy until there is a request in transit
		Busy: make(map[string]bool),

		// Map the reference code to the response from the target forwarded by the agent
		ReferenceToResponse: make(map[string]chan *http.Request),

		// Inform the server that the response was sent to the external requester so the
		// connection with the agent can be closed.
		responseSentChannel: make(map[string]chan struct{}),
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
	if action == "pull" {
		l.handlePullRequest(code, w, r)
		return
	} else if action == "forward" {
		l.handleForwardRequest(code, w, r)
		return
	} else if action == "forward_response" {
		l.handleForwardResponse(code, w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "wrong request")

}
