package biz

import (
	"context"
	"github.com/gorilla/websocket"
	"undersea/im_balance/conf"
	"undersea/pkg/log"
	"undersea/pkg/message"
)

type BalanceUseCase struct {
	conf conf.Conf
}

func NewBalanceUseCase(conf conf.Conf) *BalanceUseCase {
	return &BalanceUseCase{conf: conf}
}

func (uc *BalanceUseCase) HandlePickIpMessage(ctx context.Context, conn *websocket.Conn) (err error) {
	var ip string
	// 这里使用随机算法获取ip
	imServer.ipMap.Range(func(key, value any) bool {
		ip = key.(string)
		return false
	})

	if ip == "" {
		log.E(ctx, err).Msgf("所有im节点均不可用")
		return
	}

	// 发送消息给客户端
	err = message.SendWebSocketMessage(ctx, conn, &message.PickNodeIpReplyMessage{
		Ip: ip,
	}, message.MesTypeReplyPickNodeIp, "")
	if err != nil {
		log.E(ctx, err).Msgf("所有im节点均不可用")
		return
	}
	return
}
