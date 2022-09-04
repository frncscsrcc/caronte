package agent

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	ServerSecret   string
	ServerHost     string
	ServerPort     string
	ServerProtocol string
	Timeout        int
	TargetHost     string
	TargetPort     string
	TargetProtocol string
	Code           string
}

var agentConfig Config

func Run(config Config) {
	agentConfig = config

	serverURL := fmt.Sprintf("%s://%s:%s", agentConfig.ServerProtocol, agentConfig.ServerHost, agentConfig.ServerPort)

	log.Printf("Agent connecting to %s", serverURL)
	for {
		req, _ := http.NewRequest("GET", serverURL+"/receive/"+agentConfig.Code, new(bytes.Buffer))

		req.Header.Set("X-TIMEOUT", fmt.Sprintf("%d", agentConfig.Timeout))
		req.Header.Set("Authorization", "Bearer: "+agentConfig.ServerSecret)

		serverResponse, err := http.DefaultClient.Do(req)

		if serverResponse.StatusCode == http.StatusRequestTimeout {
			continue
		}

		if err != nil {
			log.Print("ERROR: " + err.Error())
		}

		if err != nil || serverResponse.StatusCode != http.StatusOK {
			time.Sleep(5 * time.Second)
			log.Printf("Agent connecting to %s", serverURL)
		} else {
			handleServerResponse(serverResponse)
		}
	}

}

func handleServerResponse(r *http.Response) error {
	path := r.Header.Get("X-ORIGINAL-PATH")
	method := r.Header.Get("X-ORIGINAL-METHOD")

	url := fmt.Sprintf("%s://%s:%s",
		agentConfig.TargetProtocol,
		agentConfig.TargetHost,
		agentConfig.TargetPort,
	) + path

	log.Printf("[%s] %s", method, url)

	req, err := http.NewRequest(strings.ToUpper(method), url, r.Body)
	if err != nil {
		return err
	}
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Set(name, value)
		}
	}

	log.Printf("Forwarding request to %s\n", url)

	_, errReq := http.DefaultClient.Do(req)
	return errReq
}
