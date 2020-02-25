package agentcmd

import "github.com/spf13/cobra"

var Agent = &cobra.Command{
	Use:   "agent",
	Short: "Agent commands",
}

func init() {
	Agent.AddCommand(start)
	Agent.AddCommand(initAgentCMD)
}
