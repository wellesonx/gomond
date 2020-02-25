package servercmd

import (
	"encoding/json"
	"github.com/gelleson/gomond/gomond/server"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var start = &cobra.Command{
	Use:   "start",
	Short: "start server",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.Flags().StringP("config", "c", "server.json", "--config server.json")

	},

	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal(err)
		}

		file, err := ioutil.ReadFile(config)
		if err != nil {
			log.Fatal(err)
		}

		option := server.Option{}

		err = json.Unmarshal(file, &option)

		if err != nil {
			log.Fatal(err)
		}

		newServer, err := server.NewServer(option)
		if err != nil {
			log.Fatal(err)
		}

		newServer.Run()
	},
}
