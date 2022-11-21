/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"webhookme/pkg/server"

	"github.com/spf13/cobra"
)

var host string
var port string
var maxTimeout int
var bufferSize int
var secret string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Run(server.Config{
			Host:       host,
			Port:       port,
			MaxTimeout: maxTimeout,
			BufferSize: bufferSize,
			Secret:     secret,
		})
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&host, "host", "0.0.0.0", "listening host")
	serverCmd.Flags().StringVar(&port, "port", "8080", "listening port")
	serverCmd.Flags().IntVar(&maxTimeout, "max-timeout", 120, "max timeout allowed server side")
	serverCmd.Flags().IntVar(&bufferSize, "buffer", 5, "max number of buffered requests, per resource")
	serverCmd.Flags().StringVar(&secret, "secret", "", "Instance secret")
}
