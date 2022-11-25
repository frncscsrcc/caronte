/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"caronte/pkg/agent"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var proxyProtocol string
var proxyHost string
var proxyPort string

var targetProtocol string
var targetHost string
var targetPort string
var targetSendReply bool

var agentCode string
var agentSecret string
var agentTimeout int

// agentCmd represents the client command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run an agent to expose an internal service (target), to an external caller, via an external proxy",

	Run: func(cmd *cobra.Command, args []string) {
		if agentCode == "" {
			fmt.Print("Missing --agent-code")
			return
		}

		agent.Run(agent.Config{
			ProxyProtocol:   proxyProtocol,
			ProxyHost:       proxyHost,
			ProxyPort:       proxyPort,
			Secret:          agentSecret,
			TargetProtocol:  targetProtocol,
			TargetHost:      "localhost",
			TargetPort:      targetPort,
			TargetSendReply: targetSendReply,
			AgentCode:       agentCode,
			Timeout:         agentTimeout,
		})
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringVar(&proxyProtocol, "proxy-protocol", "http", "Server protocol")
	agentCmd.Flags().StringVar(&proxyHost, "proxy-host", "localhost", "Server host")
	agentCmd.Flags().StringVar(&proxyPort, "proxy-port", "8080", "Server port")

	agentCmd.Flags().StringVar(&targetProtocol, "target-protocol", "http", "Server protocol")
	//agentCmd.Flags().StringVar(&targetHost, "target-host", "localhost", "Server host")
	agentCmd.Flags().StringVar(&targetPort, "target-port", "8080", "Server port")
	agentCmd.Flags().BoolVar(&targetSendReply, "send-reply", false, "Allow the agent to send the replies to the caller")

	agentCmd.Flags().StringVar(&agentCode, "agent-code", "", "Agent uniq resource identifier")
	agentCmd.Flags().StringVar(&agentSecret, "secret", os.Getenv("CARONTE_SECRET"), "Shared secret. If not passed the ENV CARONTE_SECRET will be used instead.")
	agentCmd.Flags().IntVar(&agentTimeout, "timeout", 60, "Requested timeout")
}
