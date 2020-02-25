package cmd

import (
	"github.com/gelleson/gomond/gomond/cmd/agent"
	"github.com/gelleson/gomond/gomond/cmd/servercmd"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use: "gomond",
}

func init() {

	root.AddCommand(agent.Agent)
	root.AddCommand(servercmd.CMD)
}

func Execute() error {
	return root.Execute()
}
