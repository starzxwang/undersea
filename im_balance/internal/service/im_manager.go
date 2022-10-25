package service

import (
	"context"
	"net"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/biz"
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

func (s *ImManagerService) ProcessImManager(conn net.Conn) {
	ctx := context.Background()
	defer conn.Close()
	for {

		buf := make([]byte, 1024)

		log.I(ctx).Msgf("服务器在等待客户端发送信息,%s", conn.RemoteAddr().String())
		n, err := conn.Read(buf)
		if err != nil {
			log.E(ctx, err).Msgf("服务器端的read err=" + err.Error())
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
		log.E(ctx, err).Msgf("ConvertBytes2Message err")
		return
	}

	switch mes.Type {
	case message.MesTypeHeartBeat:
		// 心跳
		err = s.imManagerUseCase.HandleHeartBeatMessage(ctx, conn, mes)
	default:
	}
}
