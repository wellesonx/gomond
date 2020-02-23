package cmd

import (
	"github.com/gelleson/gomond/gomond/cmd/agent"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use: "gomond",
}

func init() {

	root.AddCommand(agent.Agent)
}

func Execute() error {
	return root.Execute()
}
