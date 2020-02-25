package cmd

import (
	"github.com/gelleson/gomond/gomond/cmd/agentcmd"
	"github.com/gelleson/gomond/gomond/cmd/servercmd"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "gomond",
	Short: "monitoring system",
}

func init() {

	root.AddCommand(agentcmd.Agent)
	root.AddCommand(servercmd.CMD)
}

func Execute() error {
	return root.Execute()
}
