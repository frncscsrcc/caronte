/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"caronte/pkg/proxy"
	"os"

	"github.com/spf13/cobra"
)

var host string
var port string
var maxTimeoutLongPoll int
var maxTimeoutResponse int
var bufferSize int
var proxySecret string

// proxyCmd represents the server command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Run a proxy to forward the connect an external caller to an internal service (target)",
	Run: func(cmd *cobra.Command, args []string) {
		proxy.Run(proxy.Config{
			Host:               host,
			Port:               port,
			MaxTimeoutLongPoll: maxTimeoutLongPoll,
			MaxTimeoutResponse: maxTimeoutResponse,
			BufferSize:         bufferSize,
			Secret:             proxySecret,
		})
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().StringVar(&host, "host", "0.0.0.0", "listening host")
	proxyCmd.Flags().StringVar(&port, "port", "8080", "listening port")
	proxyCmd.Flags().IntVar(&maxTimeoutLongPoll, "max-timeout-long-poll", 120, "max timeout long poll")
	proxyCmd.Flags().IntVar(&maxTimeoutResponse, "max-timeout-response", 120, "max timeout to get a response from the agent")
	proxyCmd.Flags().IntVar(&bufferSize, "buffer", 5, "max number of buffered requests, per resource")
	proxyCmd.Flags().StringVar(&proxySecret, "secret", os.Getenv("CARONTE_SECRET"), "Shared secret. If not passed the ENV CARONTE_SECRET will be used instead.")
}
