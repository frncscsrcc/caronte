package client

import (
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

	Timeout int

	TargetHost     string
	TargetPort     string
	TargetProtocol string

	Code string
}

var clientConfig Config

func Run(config Config) {
	clientConfig = config

	serverURL := fmt.Sprintf("%s://%s:%s", clientConfig.ServerProtocol, clientConfig.ServerHost, clientConfig.ServerPort)
	log.Printf("Client connecting to %s", serverURL)

	for {
		serverResponse, err := http.Get(serverURL + "/receive/" + clientConfig.Code)

		if serverResponse.StatusCode == http.StatusRequestTimeout {
			continue
		}

		if err != nil {
			log.Print("ERROR: " + err.Error())
			time.Sleep(5 * time.Second)
		} else {
			handleServerResponse(serverResponse)
		}
	}

}

func handleServerResponse(r *http.Response) error {
	path := r.Header.Get("X-ORIGINAL-PATH")
	method := r.Header.Get("X-ORIGINAL-METHOD")
	// for name, values := range r.Header {

	// }
	url := fmt.Sprintf("%s://%s:%s",
		clientConfig.TargetProtocol,
		clientConfig.TargetHost,
		clientConfig.TargetPort,
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

	targetResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	fmt.Print(targetResponse)
	return nil
}
