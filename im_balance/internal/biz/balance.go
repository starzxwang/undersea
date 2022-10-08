package biz

import (
	"context"
	"github.com/gorilla/websocket"
	"undersea/im_balance/conf"
	error_v2 "undersea/pkg/err"
	"undersea/pkg/message"
)



type BalanceUseCase struct {
	conf conf.Conf
}


func NewBalanceUseCase(conf conf.Conf) *BalanceUseCase {
	return &BalanceUseCase{conf: conf}
}

func (uc *BalanceUseCase) HandlePickNodeIpMessage(ctx context.Context, conn *websocket.Conn) (err error) {
	var ip string
	imServer.rwLock.RLock()
	ipList := imServer.ipList
	imServer.rwLock.RUnlock()
	for currentIp := range ipList {
		// 由于map的key遍历时有随机性，我们直接取第一个作为返回ip
		ip = currentIp
		break
	}

	if ip == "" {
		err = error_v2.PrintError(ctx, err, "所有im节点均不可用")
		return
	}

	// 发送消息给客户端
	err = message.SendWebSocketMessage(ctx, conn, &message.PickNodeIpReplyMessage{
		Ip: ip,
	}, message.MesTypeReplyPickNodeIp, "")
	if err != nil {
		err = error_v2.PrintError(ctx, err, "所有im节点均不可用")
		return
	}
	return
}