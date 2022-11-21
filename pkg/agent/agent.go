package agent

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerSecret    string
	ServerHost      string
	ServerPort      string
	ServerProtocol  string
	Timeout         int
	TargetHost      string
	TargetPort      string
	TargetProtocol  string
	TargetSendReply bool
	Code            string
}

var agentConfig Config

func getServerURL() string {
	return fmt.Sprintf("%s://%s:%s", agentConfig.ServerProtocol, agentConfig.ServerHost, agentConfig.ServerPort)
}

func Run(config Config) {
	agentConfig = config

	serverURL := getServerURL()

	log.Printf("Agent connecting to %s", serverURL)
	for {
		req, _ := http.NewRequest("GET", serverURL+"/pull/"+agentConfig.Code, new(bytes.Buffer))

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
			log.Printf("Agent connecting to %s [status=%d]", serverURL, serverResponse.StatusCode)
		} else {
			handleServerResponse(serverResponse)
		}
	}

}

func handleServerResponse(r *http.Response) error {
	path := r.Header.Get("X-ORIGINAL-PATH")
	method := r.Header.Get("X-ORIGINAL-METHOD")
	responseReference := r.Header.Get("X-RESPONSE-REFERENCE")

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

	response, errReq := http.DefaultClient.Do(req)
	if errReq != nil {
		return errReq
	}

	if agentConfig.TargetSendReply {
		return forwardResponse(responseReference, response)
	} else {
		return forwardAck(responseReference)
	}

}

func forwardAck(responceReference string) error {
	req, _ := http.NewRequest(
		"POST",
		getServerURL()+"/forward_response/"+agentConfig.Code,
		bytes.NewBuffer([]byte{}),
	)
	req.Header.Set("X-RESPONSE-REFERENCE", responceReference)
	req.Header.Set("X-STATUS", "200")
	_, err := http.DefaultClient.Do(req)
	return err
}

func forwardResponse(responceReference string, response *http.Response) error {
	log.Printf("Sending reply to %s", getServerURL()+"/forward_response/"+agentConfig.Code)

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	req, _ := http.NewRequest(
		"POST",
		getServerURL()+"/forward_response/"+agentConfig.Code,
		bytes.NewBuffer(body),
	)
	req.Header.Set("X-RESPONSE-REFERENCE", responceReference)
	req.Header.Set("X-STATUS", strconv.Itoa(response.StatusCode))
	for name, values := range response.Header {
		for _, value := range values {
			req.Header.Set(name, value)
		}
	}
	serverResponse, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	Skip(serverResponse)
	return nil
}

func Skip(x interface{}) {}
