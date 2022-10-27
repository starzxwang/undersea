package server

import (
	"context"
	"google.golang.org/grpc"
	"time"
	v1 "undersea/api/im-balance/v1"
	"undersea/im-manage/conf"
	"undersea/pkg/log"
	"undersea/pkg/util"
)

const (
	ServiceName = "im-manage"
)

type GrpcClient struct {
	conf                conf.Conf
	conn                *grpc.ClientConn
	imBalanceGrpcClient v1.HeartBeatServiceClient
}

func NewGrpcClient(ctx context.Context, conf conf.Conf) *GrpcClient {
	return &GrpcClient{
		conf: conf,
	}
}

func (g *GrpcClient) Name() string {
	return "grpc"
}

func (g *GrpcClient) Start(ctx context.Context) error {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(g.conf.Grpc.ImBalanceAddr, opts...)
	if err != nil {
		log.E(ctx, err).Msgf("grpc.Dial err")
		return err
	}

	g.imBalanceGrpcClient = v1.NewHeartBeatServiceClient(conn)
	go g.balanceHeartBeat()
	log.I(ctx).Msgf("[%s] client start %s", g.Name(), g.conf.Grpc.ImBalanceAddr)
	return nil
}

func (g *GrpcClient) balanceHeartBeat() {
	ticker := time.Tick(g.conf.Ws.HeartBeatInterval)
	ctx := context.Background()
	for range ticker {
		_, err := g.imBalanceGrpcClient.Ping(ctx, &v1.PingReq{
			ServiceName: ServiceName,
			Ip:          util.GetIpAddr(),
		})

		if err != nil {
			log.E(ctx, err).Msgf("balanceHeartBeat->ping err")
			time.Sleep(time.Second)
			continue
		}

	}
}

func (g *GrpcClient) Stop(ctx context.Context) error {
	g.conn.Close()
	log.I(ctx).Msgf("[%s] client stopping", g.Name())
	return nil
}
