/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var echoHost, echoPort string

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
}

func handleRequest() {
	http.HandleFunc("/", ServeHTTP)
	hostAndPort := echoHost + ":" + echoPort
	log.Printf("Listening on %s", hostAndPort)
	if err := http.ListenAndServe(hostAndPort, nil); err != nil {
		panic(err)
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("New Request")
	log.Print("Method: " + r.Method)
	log.Print("Path: " + r.URL.Path)
	log.Print("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("\t%s: %s\n", name, value)
		}
	}
	buf := new(strings.Builder)
	io.Copy(buf, r.Body)
	r.Body.Close()
	log.Print("Body:")
	log.Print("\t", buf.String())
	log.Print("-------------")
}
