/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var echoHost, echoPort string
var delay int

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleRequest()
	},
}

func init() {
	rootCmd.AddCommand(echoCmd)
	echoCmd.Flags().StringVar(&echoHost, "host", "0.0.0.0", "listening host")
	echoCmd.Flags().StringVar(&echoPort, "port", "5000", "listening port")
	echoCmd.Flags().IntVar(&delay, "delay", 1, "delay in sec")
}

func handleRequest() {
	http.HandleFunc("/", ServeHTTP)
	hostAndPort := echoHost + ":" + echoPort
	log.Printf("Listening on %s (response delay set to %d sec)", hostAndPort, delay)
	if err := http.ListenAndServe(hostAndPort, nil); err != nil {
		panic(err)
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if delay > 0 {
		log.Printf("Received request, but waiting %d sec(s).\n", delay)
		time.Sleep(time.Second * time.Duration(delay))
	}

	log.Print("New Request")
	log.Print("Method: " + r.Method)
	log.Print("Path: " + r.URL.Path)
	log.Print("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("\t%s: %s\n", name, value)
			w.Header().Add(name, value)
		}
	}
	buf := new(strings.Builder)
	io.Copy(buf, r.Body)
	log.Print("Body:")
	log.Print("\t", buf.String())
	log.Print("-------------")
	fmt.Fprint(w, buf)
	r.Body.Close()
}
