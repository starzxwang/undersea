package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "undersea/api/im-balance/v1"
	"undersea/im-balance/conf"
	"undersea/im-balance/internal/biz"
	"undersea/pkg/message"
)

type HeartBeatService struct {
	v1.UnimplementedHeartBeatServiceServer
	conf             conf.Conf
	heartBeatUseCase *biz.HeartBeatUseCase
}

func NewHeartBeatService(conf conf.Conf, heartBeatUseCase *biz.HeartBeatUseCase) v1.HeartBeatServiceServer {
	return &HeartBeatService{
		conf:             conf,
		heartBeatUseCase: heartBeatUseCase,
	}
}

func (s *HeartBeatService) Ping(ctx context.Context, req *v1.PingReq) (resp *v1.PingResp, err error) {
	if req.ServiceName == "" || req.Ip == "" {
		return nil, status.Errorf(codes.InvalidArgument, "参数不能为空")
	}

	s.heartBeatUseCase.SaveHeartBeat(s.convertReq2HeartBeatMessage(req))
	return &v1.PingResp{}, nil
}

func (s *HeartBeatService) convertReq2HeartBeatMessage(req *v1.PingReq) (ret *message.HeartBeatMessage) {
	return &message.HeartBeatMessage{
		ServiceName: req.ServiceName,
		Ip:          req.Ip,
	}
}
