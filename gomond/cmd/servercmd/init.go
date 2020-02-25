package servercmd

import "github.com/spf13/cobra"

var CMD = &cobra.Command{
	Use: "server",
}

func init() {
	CMD.AddCommand(start)
}
