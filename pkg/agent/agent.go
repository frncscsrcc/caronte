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
	Secret          string
	ProxyHost       string
	ProxyPort       string
	ProxyProtocol   string
	Timeout         int
	TargetHost      string
	TargetPort      string
	TargetProtocol  string
	TargetSendReply bool
	AgentCode       string
}

var agentConfig Config

func getProxyURL() string {
	return fmt.Sprintf("%s://%s:%s", agentConfig.ProxyProtocol, agentConfig.ProxyHost, agentConfig.ProxyPort)
}

func Run(config Config) {
	agentConfig = config

	serverURL := getProxyURL()

	log.Printf("Agent connecting to the proxy %s", serverURL)
	for {
		req, _ := http.NewRequest("GET", serverURL+"/pull/"+agentConfig.AgentCode, new(bytes.Buffer))

		req.Header.Set("X-TIMEOUT", fmt.Sprintf("%d", agentConfig.Timeout))
		req.Header.Set("Authorization", "Bearer: "+agentConfig.Secret)

		proxyResponse, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Printf("[ERROR] %s. Trying again in 5 sec...\n", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		if proxyResponse.StatusCode == http.StatusRequestTimeout {
			continue
		}

		if err != nil {
			log.Print("ERROR: " + err.Error())
		}

		if err != nil || proxyResponse.StatusCode != http.StatusOK {
			log.Printf("Proxy responded %d, trying again in 5 sec...", proxyResponse.StatusCode)
			time.Sleep(5 * time.Second)
		} else {
			handleServerResponse(proxyResponse)
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

	log.Printf("Forwarding request to the target %s\n", url)

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
		getProxyURL()+"/forward_response/"+agentConfig.AgentCode,
		bytes.NewBuffer([]byte{}),
	)
	req.Header.Set("X-RESPONSE-REFERENCE", responceReference)
	req.Header.Set("X-STATUS", "200")
	_, err := http.DefaultClient.Do(req)
	return err
}

func forwardResponse(responceReference string, response *http.Response) error {
	log.Printf("Forwarding the response to the proxy %s", getProxyURL()+"/forward_response/"+agentConfig.AgentCode)

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	req, _ := http.NewRequest(
		"POST",
		getProxyURL()+"/forward_response/"+agentConfig.AgentCode,
		bytes.NewBuffer(body),
	)
	req.Header.Set("X-RESPONSE-REFERENCE", responceReference)
	req.Header.Set("X-STATUS", strconv.Itoa(response.StatusCode))
	for name, values := range response.Header {
		for _, value := range values {
			req.Header.Set(name, value)
		}
	}
	_, sendError := http.DefaultClient.Do(req)

	if sendError != nil {
		return err
	}
	return nil
}
