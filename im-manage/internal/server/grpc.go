package server

import (
	"context"
	"google.golang.org/grpc"
	"time"
	v1 "undersea/api/im-balance/v1"
	"undersea/im-manage/conf"
	"undersea/pkg/log"
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
	ticker := time.NewTicker(g.conf.Ws.HeartBeatInterval)
	ctx := context.Background()
	isInit := true
	for {
		select {
		case <-ticker.C:
			_, err := g.imBalanceGrpcClient.Ping(ctx, &v1.PingReq{
				ServiceName: ServiceName,
				Ip:          g.getWsIp(),
			})
			if err != nil {
				log.E(ctx, err).Msgf("balanceHeartBeat->ping err")
				time.Sleep(time.Second * 2)
			}
		default:
			if isInit {
				_, err := g.imBalanceGrpcClient.Ping(ctx, &v1.PingReq{
					ServiceName: ServiceName,
					Ip:          g.getWsIp(),
				})
				if err != nil {
					log.E(ctx, err).Msgf("balanceHeartBeat->ping err")
					time.Sleep(time.Second * 2)
				}

				isInit = false
			}

		}
	}
}

func (g *GrpcClient) Stop(ctx context.Context) error {
	g.conn.Close()
	log.I(ctx).Msgf("[%s] client stopping", g.Name())
	return nil
}

func (g *GrpcClient) getWsIp() string {
	return "ws://" + g.conf.Ws.Addr + "/ws"
}
