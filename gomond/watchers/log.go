package watchers

import (
	"github.com/gelleson/gomond/gomond/collector"
	"github.com/gelleson/gomond/gomond/parser"
	"github.com/gelleson/gomond/gomond/provider"
	"github.com/sirupsen/logrus"
)

type Log struct {
	provider     provider.Provider
	parser       parser.Parser
	logCollector collector.LogCollector
	logger       *logrus.Logger
	messagesCh   chan []byte
}

func (l *Log) SetParser(parser parser.Parser) {
	l.parser = parser
}

func (l *Log) SetLogCollector(logCollector collector.LogCollector) {
	l.logCollector = logCollector
}

func (l *Log) SetProvider(provider provider.Provider) {
	l.provider = provider
}

func NewLogApp(provider provider.Provider, parser parser.Parser, logCollector collector.LogCollector, logger *logrus.Logger) *Log {
	return &Log{provider: provider, parser: parser, logCollector: logCollector, logger: logger, messagesCh: make(chan []byte)}
}

func (l *Log) Run() {

	err := l.provider.Start()

	if err != nil {
		l.logger.Error(err)
		return
	}

	go l.startWorker()

	err = l.provider.Follow(l.messagesCh)

	if err != nil {
		l.logger.Error(err)
	}

	err = l.provider.Close()

	if err != nil {
		l.logger.Error(err)
	}
}

func (l Log) startWorker() {
	for message := range l.messagesCh {
		parsedLog, err := l.parser.Parse(message)
		if err != nil {
			l.logger.Error(err)
		}
		l.logCollector.Push(parsedLog)
	}
}
