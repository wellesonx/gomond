package agent

import "github.com/spf13/cobra"

var Agent = &cobra.Command{
	Use: "agent",
}

func init() {
	Agent.AddCommand(start)
}
