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
var targetMaxAttempts int

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
			ProxyProtocol:     proxyProtocol,
			ProxyHost:         proxyHost,
			ProxyPort:         proxyPort,
			Secret:            agentSecret,
			TargetProtocol:    targetProtocol,
			TargetHost:        targetHost,
			TargetPort:        targetPort,
			TargetSendReply:   targetSendReply,
			TargetMaxAttempts: targetMaxAttempts,
			AgentCode:         agentCode,
			Timeout:           agentTimeout,
		})
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringVar(&proxyProtocol, "proxy-protocol", "http", "Proxy protocol")
	agentCmd.Flags().StringVar(&proxyHost, "proxy-host", "localhost", "Proxy host")
	agentCmd.Flags().StringVar(&proxyPort, "proxy-port", "8080", "Proxy port")

	agentCmd.Flags().StringVar(&targetProtocol, "target-protocol", "http", "Target protocol")
	agentCmd.Flags().StringVar(&targetHost, "target-host", "localhost", "Target host")
	agentCmd.Flags().StringVar(&targetPort, "target-port", "8080", "Target port")
	agentCmd.Flags().IntVar(&targetMaxAttempts, "target-max-attempts", 5, "How many times the agent will try to call the target in case of errors")

	agentCmd.Flags().BoolVar(&targetSendReply, "send-reply", false, "Allow the agent to send the replies to the caller")

	agentCmd.Flags().StringVar(&agentCode, "agent-code", "", "Agent uniq resource identifier")
	agentCmd.Flags().StringVar(&agentSecret, "secret", os.Getenv("CARONTE_SECRET"), "Shared secret. If not passed the ENV CARONTE_SECRET will be used instead.")
	agentCmd.Flags().IntVar(&agentTimeout, "timeout", 60, "Requested timeout")
}
