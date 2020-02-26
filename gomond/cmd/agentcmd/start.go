package agentcmd

import (
	"encoding/json"
	"github.com/gelleson/gomond/gomond/agent"
	"github.com/gelleson/gomond/gomond/collector"
	"github.com/gelleson/gomond/gomond/factory"
	"github.com/gelleson/gomond/gomond/watchers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Watcher struct {
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Parser   string           `json:"parser"`
	Provider *json.RawMessage `json:"provider"`
	Option   *json.RawMessage `json:"option"`
}

var start = &cobra.Command{
	Use:   "start",
	Short: "Start agent server",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Duration("ttl", time.Hour, "--ttl 12h")

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

		configRaw := make(map[string]*json.RawMessage)

		err = json.Unmarshal(file, &configRaw)

		if err != nil {
			log.Fatal(err)
		}

		watchersConfig := make(map[string]*json.RawMessage)

		err = json.Unmarshal(*configRaw["watchers"], &watchersConfig)

		if err != nil {
			log.Fatal(err)
		}

		logConfig := agent.Log{}

		err = json.Unmarshal(*configRaw["log"], &logConfig)

		if err != nil {
			log.Fatal(err)
		}

		logger := logrus.New()

		logger.SetReportCaller(true)

		logger.SetLevel(logConfig.Level)
		logger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: false,
		})

		output, err := os.OpenFile(logConfig.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0744)
		if err != nil {
			logger.Fatal(err)
		}

		logger.SetOutput(output)

		duration, err := cmd.Flags().GetDuration("ttl")

		if err != nil {
			logger.Fatal(err)
		}

		collection := collector.NewMemoryLogCollector(collector.MemoryOption{
			TTL: duration,
		})

		logger.Info("inited")

		agent := agent.NewAgent(agent.GRPC{Port: 2001}, collection, logger)

		for key, message := range watchersConfig {

			data := make([]Watcher, 0)

			err := json.Unmarshal(*message, &data)

			if err != nil {
				logger.Fatal(err)
			}
			for _, watcher := range data {
				parser, err := factory.Parser(key, watcher.Name, watcher.Parser, watcher.Option)
				if err != nil {
					log.Fatal(err)
				}

				provider, err := factory.Provider("file", watcher.Provider)
				if err != nil {
					logger.Fatal(err)
				}

				watcher := watchers.NewLogApp(provider, parser, collection, logger)

				agent.AddWatcher(watcher)
			}

		}

		agent.StartWatchers()
	},
}

func init() {
	start.Flags().StringP("config", "c", "config.json", "-c config.json")
}
