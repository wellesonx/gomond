package agentcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var initAgentCMD = &cobra.Command{
	Use:   "init",
	Short: "Init config file for agent",
	PreRun: func(cmd *cobra.Command, args []string) {

		cmd.Flags().StringP("name", "n", "agent.json", "--name agent.json")

	},
	Run: func(cmd *cobra.Command, args []string) {
		config := `{
  "notification": {
	"enable": false,
	"token": ""
  },
  "log": {
	"path": "app.log",
	"level": "info"
  },
  "grpc": {
	"port": 2001
  },
  "watchers": {
	"app": [
	  {
		"name": "app-log-1",
		"type": "file",
		"parser": "json",
		"provider": {
		  "height": 1,
		  "file_name": "app.log"
		},
		"option": {
		  "message": "msg",
		  "timestamp": "time",
		  "level": "level",
		  "file": "file",
		  "line": "func"
		}
	  }
	]
  }
}`

		name, err := cmd.Flags().GetString("name")

		if err != nil {
			log.Fatal(err)
		}

		file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0744)

		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(config)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("DONE!")

	},
}
