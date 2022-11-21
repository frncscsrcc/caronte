package server

import (
	"math/rand" // todo: use crypto library
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getTimeout(d time.Duration) chan struct{} {
	timeout := make(chan struct{})
	go func() {
		time.Sleep(d)
		timeout <- struct{}{}
	}()
	return timeout
}

func getOriginalPath(r *http.Request) string {
	return "/" + strings.Join(
		strings.Split(r.URL.Path, "/")[3:],
		"/",
	)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getRandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
