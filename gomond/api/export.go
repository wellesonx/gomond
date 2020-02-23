package api

import (
	"context"
	"fmt"
	"github.com/gelleson/gomond/gomond/collector"
	"github.com/gelleson/gomond/gomond/pb"
	"github.com/mitchellh/hashstructure"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type ExportAPI struct {
	pb.UnimplementedLogAPIServer
	collect collector.LogCollector
	logger  *logrus.Logger
}

func NewExportAPI(collect collector.LogCollector, logger *logrus.Logger) *ExportAPI {
	return &ExportAPI{collect: collect, logger: logger}
}

func (e *ExportAPI) Latest(ctx context.Context, req *pb.AppLatestRequest) (*pb.AppResponse, error) {
	logs := e.collect.Get()

	reduceLogs := make([]*pb.Log, len(logs))

	appName := ""

	hostname, err := os.Hostname()

	if err != nil {
		e.logger.Error(err)
		return &pb.AppResponse{}, status.Error(codes.Internal, err.Error())
	}

	for i, log := range logs {
		appName = log.App
		reduceLogs[i] = &pb.Log{
			App:      log.App,
			Level:    log.Level.ToPBLevel(),
			Payload:  log.Payload,
			Time:     log.Timestamp.UnixNano(),
			Message:  log.Message,
			File:     log.File,
			Line:     int32(log.Line),
			Type:     pb.Type_file,
			Label:    log.LogName,
			Hostname: hostname,
		}
	}

	hash, err := hashstructure.Hash(reduceLogs, nil)

	if err != nil {
		e.logger.Error(err)
		return &pb.AppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.AppResponse{
		Logs: reduceLogs,
		App:  appName,
		Hash: fmt.Sprintf("%v", hash),
	}, nil
}

func (e ExportAPI) Live(req *pb.Empty, srv pb.LogAPI_LiveServer) error {
	for {
		logs := e.collect.Get()
		if len(logs) > 0 {
			for _, log := range logs {
				hostname, err := os.Hostname()
				if err != nil {
					return status.Error(codes.Internal, err.Error())
				}
				lg := &pb.Log{
					App:      log.App,
					Level:    log.Level.ToPBLevel(),
					Payload:  log.Payload,
					Time:     log.Timestamp.UnixNano(),
					Message:  log.Message,
					File:     log.File,
					Line:     int32(log.Line),
					Type:     pb.Type_file,
					Hostname: hostname,
					Label:    log.LogName,
				}
				err = srv.Send(lg)
				if err != nil {
					return status.Error(codes.Internal, err.Error())
				}
			}
		}
		time.Sleep(time.Millisecond * 20)
	}
}
