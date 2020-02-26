package agent

import (
	"fmt"
	"github.com/gelleson/gomond/gomond/api"
	"github.com/gelleson/gomond/gomond/collector"
	"github.com/gelleson/gomond/gomond/pb"
	"github.com/gelleson/gomond/gomond/watchers"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type GRPC struct {
	Port int `json:"port"`
}

type Notification struct {
	Enable bool   `json:"enable"`
	Token  string `json:"token"`
}

type Log struct {
	Path  string       `json:"path"`
	Level logrus.Level `json:"level"`
}

type Agent struct {
	watchers []watchers.Watcher
	logger   *logrus.Logger
	exporter *grpc.Server
	option   GRPC
	collect  collector.LogCollector
}

func NewAgent(option GRPC, collection collector.LogCollector, logger *logrus.Logger) *Agent {

	server := grpc.NewServer()

	return &Agent{
		logger:   logger,
		exporter: server,
		option:   option,
		collect:  collection,
		watchers: make([]watchers.Watcher, 0),
	}
}

func (a Agent) StartWatchers() {
	for _, watcher := range a.watchers {
		go watcher.Run()
	}

	apiExporter := api.NewExportAPI(a.collect, a.logger)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.option.Port))

	if err != nil {
		a.logger.Fatal(err)
	}

	pb.RegisterLogAPIServer(a.exporter, apiExporter)

	err = a.exporter.Serve(l)

	if err != nil {
		a.logger.Fatal(err)
	}

}

func (a *Agent) AddWatcher(w watchers.Watcher) {

	a.watchers = append(a.watchers, w)
}
