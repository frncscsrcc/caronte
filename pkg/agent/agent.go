package agent

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Secret            string
	ProxyHost         string
	ProxyPort         string
	ProxyProtocol     string
	Timeout           int
	TargetHost        string
	TargetPort        string
	TargetProtocol    string
	TargetSendReply   bool
	TargetMaxAttempts int
	AgentCode         string
}

var agentConfig Config

func Run(config Config) {
	agentConfig = config

	for {
		// Connect to proxy
		proxyResponse, connectToProxyError := connectToProxy()
		if connectToProxyError != nil {
			continue
		}

		// Connect to target
		targetResponse, connectToTargetError := connectToTarget(proxyResponse)
		if connectToTargetError != nil {
			log.Printf("[ERROR] Too many attempts, giving up")
			continue
		}

		// Respond to the proxy
		responseReference := proxyResponse.Header.Get("X-RESPONSE-REFERENCE")
		respondToProxy(responseReference, targetResponse)
	}
}

func getProxyURL() string {
	return fmt.Sprintf("%s://%s:%s", agentConfig.ProxyProtocol, agentConfig.ProxyHost, agentConfig.ProxyPort)
}

func connectToProxy() (*http.Response, error) {
	proxyURL := getProxyURL()
	log.Printf("Agent connecting to the proxy %s", proxyURL)

	req, _ := http.NewRequest(
		"GET",
		proxyURL+"/pull/"+agentConfig.AgentCode,
		new(bytes.Buffer),
	)

	req.Header.Set("X-TIMEOUT", fmt.Sprintf("%d", agentConfig.Timeout))
	req.Header.Set("Authorization", "Bearer: "+agentConfig.Secret)

	proxyResponse, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("[ERROR] %s. Trying again in 5 sec...\n", err.Error())
		time.Sleep(5 * time.Second)
		return nil, err
	}

	if proxyResponse.StatusCode == http.StatusTooManyRequests {
		err := errors.New("busy")
		log.Print("[ERROR] Proxy is busy. Trying again in 5 sec...\n")
		time.Sleep(5 * time.Second)
		return nil, err
	}

	if proxyResponse.StatusCode != http.StatusOK {
		time.Sleep(5 * time.Second)
		err := errors.New(fmt.Sprintf("proxy returned status [%d]", proxyResponse.StatusCode))
		log.Print(err.Error())
		return nil, err
	}

	return proxyResponse, nil
}

func connectToTarget(proxyResponse *http.Response) (*http.Response, error) {
	for attempt := 1; attempt <= agentConfig.TargetMaxAttempts; attempt += 1 {
		targetResponse, callTargetError := handleServerResponse(proxyResponse)

		if callTargetError != nil {
			log.Printf(
				"Target returned an error (%s), trying again in 5 sec... [attempt %d of %d]",
				callTargetError.Error(),
				attempt,
				agentConfig.TargetMaxAttempts,
			)
			time.Sleep(5 * time.Second)
			continue
		}

		return targetResponse, nil
	}
	return nil, errors.New("too many attempts connecting to the target")
}

func handleServerResponse(r *http.Response) (*http.Response, error) {
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
		return nil, err
	}
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Set(name, value)
		}
	}

	response, errReq := http.DefaultClient.Do(req)
	if errReq != nil {
		return nil, errReq
	}

	return response, nil
}

func respondToProxy(responseReference string, targetResponse *http.Response) error {
	var sendResponseError error

	if agentConfig.TargetSendReply {
		sendResponseError = forwardResponse(responseReference, targetResponse)
	} else {
		sendResponseError = forwardAck(responseReference)
	}

	if sendResponseError != nil {
		log.Printf("[ERROR] Can not send the response to the proxy (%s)", sendResponseError.Error())
	}

	return sendResponseError
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

func forwardError(responceReference string) error {
	req, _ := http.NewRequest(
		"POST",
		getProxyURL()+"/forward_response/"+agentConfig.AgentCode,
		bytes.NewBuffer([]byte{}),
	)
	req.Header.Set("X-RESPONSE-REFERENCE", responceReference)
	req.Header.Set("X-STATUS", "500")
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
