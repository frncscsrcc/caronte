/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"webhookme/pkg/client"

	"github.com/spf13/cobra"
)

var clientProtocol string
var clientHost string
var clientPort string

var targetProtocol string
var targetHost string
var targetPort string

var code string
var clientSecret string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if code == "" {
			fmt.Print("Missing --code")
			return
		}

		client.Run(client.Config{
			ServerProtocol: clientProtocol,
			ServerHost:     clientHost,
			ServerPort:     clientPort,
			ServerSecret:   clientSecret,

			TargetProtocol: targetProtocol,
			TargetHost:     targetHost,
			TargetPort:     targetPort,

			Code: code,
		})
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&clientProtocol, "server-protocol", "http", "Server protocol")
	clientCmd.Flags().StringVar(&clientHost, "server-host", "localhost", "Server host")
	clientCmd.Flags().StringVar(&clientPort, "server-port", "8080", "Server port")

	clientCmd.Flags().StringVar(&targetProtocol, "target-protocol", "http", "Server protocol")
	clientCmd.Flags().StringVar(&targetHost, "target-host", "localhost", "Server host")
	clientCmd.Flags().StringVar(&targetPort, "target-port", "8080", "Server port")

	clientCmd.Flags().StringVar(&code, "code", "", "Uniq resource identifier")
	clientCmd.Flags().StringVar(&clientSecret, "secret", "", "Client secret")
}
