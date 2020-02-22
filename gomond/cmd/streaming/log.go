package main

import (
	"context"
	"github.com/gelleson/gomond/gomond/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Log struct {
	pb.UnimplementedLogAPIServer

	logs []*pb.Log
}

func NewLog() *Log {
	return &Log{
		logs: make([]*pb.Log, 0),
	}
}

func (l *Log) Append(ls pb.Log) {
	l.logs = append(l.logs, &ls)
}

func (l *Log) Flush() {
	l.logs = make([]*pb.Log, 0)
}

func (l *Log) Live(req *pb.AppLiveRequest, srv pb.LogAPI_LiveServer) error {
	return status.Errorf(codes.Unimplemented, "method Live not implemented")
}
func (l *Log) Range(ctx context.Context, req *pb.AppRangeRequest) (*pb.AppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Range not implemented")
}
func (l *Log) Latest(ctx context.Context, req *pb.AppLatestRequest) (*pb.AppResponse, error) {
	defer l.Flush()

	return &pb.AppResponse{Logs: l.logs, App: "domain", Hash: "domain"}, nil
}
