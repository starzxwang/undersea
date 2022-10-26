package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
	v1 "undersea/api/im_balance/v1"
	"undersea/im_balance/conf"
	"undersea/pkg/log"
)

type GrpcServer struct {
	*grpc.Server
	conf conf.Conf
}

func NewGrpcServer(conf conf.Conf, heartBeatServer v1.HeartBeatServiceServer) *GrpcServer {
	srv := grpc.NewServer()
	v1.RegisterHeartBeatServiceServer(srv, heartBeatServer)
	return &GrpcServer{
		Server: srv,
		conf:   conf,
	}
}

func (g *GrpcServer) Name() string {
	return "grpc"
}

func (g *GrpcServer) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", g.conf.Grpc.Addr)
	if err != nil {
		return err
	}
	log.I(ctx).Msgf("[%s] server start %s", g.Name(), g.conf.Grpc.Addr)
	return g.Serve(l)
}

func (g *GrpcServer) Stop(ctx context.Context) error {
	g.GracefulStop()
	log.I(ctx).Msgf("[%s] server stopping", g.Name())
	return nil
}
