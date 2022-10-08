package service

import (
	"context"
	"net"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/biz"
	error_v2 "undersea/pkg/err"
	"undersea/pkg/log"
	"undersea/pkg/message"
)

type ImManagerService struct {
	conf             conf.Conf
	imManagerUseCase *biz.ImManagerUseCase
}

func NewImManagerService(conf conf.Conf, imManagerUseCase *biz.ImManagerUseCase) *ImManagerService {
	return &ImManagerService{
		conf:             conf,
		imManagerUseCase: imManagerUseCase,
	}
}

// 监听TCP端口
func (s *ImManagerService) ListenTCP() {
	var err error
	ctx := context.Background()
	// 1. 监听端口
	listener, err := net.Listen("tcp", s.conf.TcpAddr)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "net listen fail")
		return
	}

	log.I(ctx).Msgf("负载均衡模块开始监听tcp端口：%s", s.conf.TcpAddr)

	defer listener.Close()

	for {
		// 2. 接收客户端请求建立链接
		conn, err := listener.Accept()
		if err != nil {
			err = error_v2.PrintError(ctx, err, "net listen accept fail")
			continue
		}

		// 3. 创建goroutine处理链接，主要是做消息转发，没有别的意义
		go s.processImManager(conn)
	}
}

func (s *ImManagerService) processImManager(conn net.Conn) {
	ctx := context.Background()
	defer conn.Close()
	for {

		buf := make([]byte, 1024)

		log.I(ctx).Msgf("服务器在等待客户端发送信息,%s", conn.RemoteAddr().String())
		n, err := conn.Read(buf)
		if err != nil {
			err = error_v2.PrintError(ctx, err, "服务器端的read err="+err.Error())
			return
			//当客户端完成任务或异常关闭后，这边我们就将协程退出，否则会循环报错链接
		}

		// 处理所有来自im_manager的消息，主要是心跳
		go s.handleImManagerMessage(conn, buf[:n])
	}
}

// 处理所有经过im_balance的TCP消息
func (s *ImManagerService) handleImManagerMessage(conn net.Conn, data []byte) {
	ctx := context.Background()
	mes, err := message.ConvertBytes2Message(ctx, data)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "ConvertBytes2Message err")
		return
	}

	switch mes.Type {
	case message.MesTypeHeartBeat:
		// 心跳
		err = s.imManagerUseCase.HandleHeartBeatMessage(ctx, conn, mes)
	default:
	}
}
