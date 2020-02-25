package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gelleson/gomond/gomond/models"
	"github.com/gelleson/gomond/gomond/pb"
	"github.com/go-macaron/binding"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/macaron.v1"
	"net/http"
	"os"
	"time"
)

type watcher struct {
	model    models.Watcher
	isActive bool
	stop     chan struct{}
	save     func(log pb.Log) error
	log      func(w watcher, message string, event models.Event)
}

func (w watcher) start(l *logrus.Logger) {
	for {
		g, err := grpc.Dial(w.model.Host, grpc.WithInsecure())
		if err != nil {
			l.Error(err)
			w.log(w, "disconnected", models.DISCONNECT)
			time.Sleep(time.Second * 15)
			continue
		}
		l.Infof("successfully connected to %s", w.model.App)
		w.log(w, "connected", models.CONNECT)

		appClient := pb.NewLogAPIClient(g)

		live, err := appClient.Live(context.Background(), &pb.Empty{})
		if err != nil {
			l.Error(err)
			w.log(w, "disconnected", models.DISCONNECT)
			time.Sleep(time.Second * 15)
			continue
		}

	Loop:
		for {
			recv, err := live.Recv()
			if err != nil {
				l.Error(err)
				w.log(w, err.Error(), models.ERROR)
				break Loop
			}

			err = w.save(*recv)
			if err != nil {
				l.Error(err)
				w.log(w, err.Error(), models.ERROR)
				break Loop
			}
		}
		time.Sleep(time.Second * 15)

	}
}

type Option struct {
	Port   int `json:"port"`
	Logger struct {
		Path  string       `json:"path"`
		Level logrus.Level `json:"level"`
	} `json:"log"`
	Database struct {
		Type string `json:"type"`
		URI  string `json:"uri"`
	}
}

type Server struct {
	option     Option
	httpServer *macaron.Macaron
	db         *gorm.DB
	watchers   []*watcher
	logger     *logrus.Logger
}

func NewServer(option Option) (*Server, error) {

	server := &Server{option: option}

	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.SetLevel(option.Logger.Level)

	f, err := os.OpenFile(option.Logger.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0744)

	if err != nil {
		return nil, err
	}
	logger.SetOutput(f)

	server.logger = logger

	server.httpServer = macaron.Classic()

	db, err := gorm.Open(option.Database.Type, option.Database.URI)
	if err != nil {
		return nil, err
	}

	server.db = db

	return server, nil
}

func (s *Server) sync() {
	watchersFromDB := make([]models.Watcher, 0)

	s.db.Where(&models.Watcher{IsActive: true}).Find(&watchersFromDB)

	watchers := make([]*watcher, len(watchersFromDB))

	for i, watcherView := range watchersFromDB {
		watchers[i] = &watcher{
			model:    watcherView,
			isActive: true,
			save: func(log pb.Log) error {

				s.db.Save(&models.Log{
					AppName:  log.App,
					Label:    log.Label,
					Hostname: log.Hostname,
					Message:  log.Message,
					Level:    int32(log.Level),
					Payload:  log.Payload,
					Line:     log.Line,
					File:     log.File,
				})

				return nil
			},
			log: func(w watcher, message string, event models.Event) {
				s.db.Save(&models.WatcherEvent{
					WatcherID: w.model.ID,
					Event:     event,
					Message:   fmt.Sprintf("%s is %s", w.model.App, message),
				})
			},
		}
	}

	s.watchers = watchers
}

func (s *Server) startWatchers() {

	for _, w := range s.watchers {
		go w.start(s.logger)
	}
}

func (s *Server) prepareServer() {
	m := s.httpServer

	s.db.AutoMigrate(&models.Log{}, &models.Watcher{}, &models.WatcherEvent{})

	s.sync()

	go s.startWatchers()

	m.Map(s.db)

	m.Group("/logs", func() {
		m.Get("/", func(dataProvider *gorm.DB, ctx *macaron.Context) (int, []byte) {
			query := dataProvider
			limit := ctx.QueryInt("limit")

			if limit == 0 {
				limit = 10
			}

			query = query.Order("created_at desc").Limit(limit)

			level := ctx.QueryInt("level")

			if level != 0 {
				query = query.Where("level >= ?", level)
			}

			app := ctx.Query("app")

			if app != "" {
				query = query.Where(&models.Log{AppName: app})
			}

			message := ctx.Query("message")

			if message != "" {
				query = query.Where("message LIKE", "%"+message+"%")
			}

			logs := make([]models.Log, 0)

			query.Find(&logs)

			logsBytes, err := json.Marshal(logs)

			if err != nil {
				return http.StatusBadGateway, []byte(err.Error())
			}

			return http.StatusOK, logsBytes
		})
	})

	m.Group("/watchers", func() {

		m.Get("/", func(dataProvider *gorm.DB) []byte {
			watchers := make([]models.Watcher, 0)
			dataProvider.Find(&watchers)
			watchersBytes, _ := json.Marshal(&watchers)
			return watchersBytes
		})

		m.Post("/", binding.Bind(models.Watcher{}), func(payload models.Watcher, dataProvider *gorm.DB) {
			dataProvider.Save(&payload)
			watcher := &watcher{
				model:    payload,
				isActive: true,
				stop:     make(chan struct{}, 0),
				save: func(log pb.Log) error {
					s.db.Save(&models.Log{
						AppName:  log.App,
						Label:    log.Label,
						Hostname: log.Hostname,
						Message:  log.Message,
						Level:    int32(log.Level),
						Payload:  log.Payload,
						Line:     log.Line,
						File:     log.File,
					})
					return nil
				},
				log: func(w watcher, message string, event models.Event) {
					s.db.Save(&models.WatcherEvent{
						WatcherID: w.model.ID,
						Event:     event,
						Message:   fmt.Sprintf("%s is %s", w.model.App, message),
					})
				},
			}

			s.watchers = append(s.watchers, watcher)

			go watcher.start(s.logger)
		})
	})
}

func (s *Server) Run() {
	s.prepareServer()
	s.httpServer.Run(s.option.Port)
}
