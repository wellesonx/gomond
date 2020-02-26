package collector

import (
	"github.com/gelleson/gomond/gomond/parser"
	"sync"
	"sync/atomic"
	"time"
)

type LogCollector interface {
	Push(log parser.Log)
	Get() []parser.Log
	Live() <-chan parser.Log
}

type MemoryOption struct {
	TTL        time.Duration `json:"ttl"`
	MaxObjects int           `json:"max_objects"`
}

type MemoryLogCollector struct {
	logs   []logMeta
	stream chan parser.Log
	viewed int64
	option MemoryOption
	mutex  sync.Mutex
}

func NewMemoryLogCollector(option MemoryOption) *MemoryLogCollector {
	collector := &MemoryLogCollector{
		logs:   make([]logMeta, 0),
		option: option,
	}

	go collector.runTTLWatcher()

	return collector
}

func (m *MemoryLogCollector) runTTLWatcher() {
	ticker := time.NewTicker(m.option.TTL)

	for _ = range ticker.C {
		m.filterByTTL()
	}
}

func (m *MemoryLogCollector) filterByTTL() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	activeLogs := make([]logMeta, 0)
	t := time.Now().Add(-m.option.TTL)

	for _, logMeta := range m.logs {

		if isActive(t, logMeta) {
			activeLogs = append(activeLogs, logMeta)
		}
	}

	m.logs = activeLogs
	m.viewed = 0
}

func (m *MemoryLogCollector) Push(log parser.Log) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logs = append(m.logs, newLogMeta(log.App, log))
}

func (m *MemoryLogCollector) Get() []parser.Log {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	logs := make([]parser.Log, 0)

	for i := m.viewed; i < int64(len(m.logs)); i++ {
		logObj := m.logs[i]
		logs = append(logs, logObj.log)
		atomic.AddInt64(&m.viewed, 1)
	}

	return logs
}

func (m *MemoryLogCollector) Live() <-chan parser.Log {
	return m.stream
}

type logMeta struct {
	app     string
	log     parser.Log
	created time.Time
}

func newLogMeta(app string, log parser.Log) logMeta {
	return logMeta{app: app, log: log, created: time.Now()}
}

func isActive(ttl time.Time, meta logMeta) bool {

	return ttl.Before(meta.created)
}
