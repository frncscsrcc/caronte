/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"webhookme/pkg/agent"

	"github.com/spf13/cobra"
)

var serverProtocol string
var serverHost string
var serverPort string

var targetProtocol string
var targetHost string
var targetPort string
var targetSendReply bool

var code string
var clientSecret string
var agentTimeout int

// agentCmd represents the client command
var agentCmd = &cobra.Command{
	Use:   "agent",
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

		agent.Run(agent.Config{
			ServerProtocol:  serverProtocol,
			ServerHost:      serverHost,
			ServerPort:      serverPort,
			ServerSecret:    clientSecret,
			TargetProtocol:  targetProtocol,
			TargetHost:      "localhost",
			TargetPort:      targetPort,
			TargetSendReply: targetSendReply,
			Code:            code,
			Timeout:         agentTimeout,
		})
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringVar(&serverProtocol, "server-protocol", "http", "Server protocol")
	agentCmd.Flags().StringVar(&serverHost, "server-host", "localhost", "Server host")
	agentCmd.Flags().StringVar(&serverPort, "server-port", "8080", "Server port")

	agentCmd.Flags().StringVar(&targetProtocol, "target-protocol", "http", "Server protocol")
	//agentCmd.Flags().StringVar(&targetHost, "target-host", "localhost", "Server host")
	agentCmd.Flags().StringVar(&targetPort, "target-port", "8080", "Server port")
	agentCmd.Flags().BoolVar(&targetSendReply, "send-reply", false, "Allow the agent to send the replies to the caller")

	agentCmd.Flags().StringVar(&code, "code", "", "Uniq resource identifier")
	agentCmd.Flags().StringVar(&clientSecret, "secret", "", "Client secret")
	agentCmd.Flags().IntVar(&agentTimeout, "timeout", 60, "Requested timeout")
}
