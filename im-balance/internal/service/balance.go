package service

import (
	"context"
	"github.com/gorilla/websocket"
	"undersea/im-balance/internal/biz"
	"undersea/pkg/log"
	"undersea/pkg/message"
)

type BalanceService struct {
	ctx            context.Context
	balanceUseCase *biz.BalanceUseCase
}

func NewBalanceService(ctx context.Context, balanceUseCase *biz.BalanceUseCase) *BalanceService {
	return &BalanceService{ctx: ctx, balanceUseCase: balanceUseCase}
}

// 这里准备根据负载均衡策略，给客户端连接分配一个im节点ip
func (s *BalanceService) HandleClientMessage(conn *websocket.Conn, data []byte) {
	var err error
	ctx := context.Background()

	// 负载均衡策略为随机分配
	// 解析消息体，获取发消息的用户uid
	mes, err := message.ConvertBytes2Message(ctx, data)

	if err != nil {
		log.E(ctx, err).Msgf("message is not valid")
		return
	}

	switch mes.Type {
	case message.MesTypePickIp:
		// 客户端需要一个im节点ip
		err = s.balanceUseCase.HandlePickIpMessage(ctx, mes, conn)
	default:
	}

	if err != nil {
		log.E(ctx, err).Msgf("handle client msg err")
		return
	}

	return
}
