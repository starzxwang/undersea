package server

import (
	"context"
	"net"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/service"
	"undersea/pkg/log"
)

type TcpServer struct {
	conf             conf.Conf
	imManagerService *service.ImManagerService
}

func NewTcpServer(conf conf.Conf, imManagerService *service.ImManagerService) *TcpServer {
	return &TcpServer{
		conf:             conf,
		imManagerService: imManagerService,
	}
}

func (*TcpServer) Name() string {
	return "tcp"
}

func (s *TcpServer) Start(ctx context.Context) (err error) {
	// 1. 监听端口
	listener, err := net.Listen("tcp", s.conf.TcpAddr)
	if err != nil {
		log.E(ctx, err).Msgf("net listen fail")
		return
	}

	log.I(ctx).Msgf("负载均衡模块开始监听tcp端口：%s", s.conf.TcpAddr)

	defer listener.Close()

	for {
		// 2. 接收客户端请求建立链接
		conn, err := listener.Accept()
		if err != nil {
			log.E(ctx, err).Msgf("net listen accept fail")
			continue
		}

		// 3. 创建goroutine处理链接，主要是做消息转发，没有别的意义
		// todo 协程池？
		go s.imManagerService.ProcessImManager(conn)
	}
}

func (*TcpServer) Stop(ctx context.Context) (err error) {
	return nil
}
