package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"undersea/im-balance/conf"
	"undersea/pkg/api"
	"undersea/pkg/log"
	"undersea/pkg/message"
)

type BalanceUseCase struct {
	conf        conf.Conf
	balanceRepo BalanceRepo
}

func NewBalanceUseCase(conf conf.Conf, balanceRepo BalanceRepo) *BalanceUseCase {
	return &BalanceUseCase{conf: conf, balanceRepo: balanceRepo}
}

func (uc *BalanceUseCase) HandlePickIpMessage(ctx context.Context, mes *message.Message, conn *websocket.Conn) (err error) {
	var sourceMes *message.PickIpMessage
	err = json.Unmarshal([]byte(mes.Data), &sourceMes)
	if err != nil {
		err = fmt.Errorf("HandlePickIpMessage->json unmarshal err,%v", err)
		return
	}

	// 取redis
	ip, err := uc.balanceRepo.GetUserIp(ctx, sourceMes.Uid)
	if err != nil {
		err = fmt.Errorf("HandlePickIpMessage->GetUserIp err,%v", err)
		message.SendExceptionWebsocketMessage(ctx, conn)
		return
	}
	if ip != "" {
		// 发送消息给客户端
		err = message.SendWebSocketMessage(ctx, conn, &message.ReplyMessageData{
			Data: &message.PickIpReplyMessage{
				Ip:  ip,
				Uid: sourceMes.Uid,
			},
		}, message.MesTypeReplyPickIp, "")
		if err != nil {
			err = fmt.Errorf("HandlePickIpMessage->SendWebSocketMessage err,%v", err)
			return
		}
	}

	// redis也没有，则重新分配ip节点
	imServer.ipMap.Range(func(key, value any) bool {
		ip = key.(string)
		return false
	})

	serviceIpList, ok := imServer.ipMap.Load(sourceMes.ServiceName)
	if !ok {
		err = message.SendWebSocketMessage(ctx, conn, &message.ReplyMessageData{
			Code:    api.CodeNotExists,
			Message: "所有im节点不可用",
		}, message.MesTypeReplyPickIp, "")
		return
	}

	serviceIpList.(*sync.Map).Range(func(key, value any) bool {
		ip = key.(string)
		return false
	})

	if ip == "" {
		err = message.SendWebSocketMessage(ctx, conn, &message.ReplyMessageData{
			Code:    api.CodeNotExists,
			Message: "所有im节点不可用",
		}, message.MesTypeReplyPickIp, "")
		return
	}

	// 将ip和user_id的映射关系放到redis，因为消息服务也会用到这个映射
	err = uc.balanceRepo.SaveIpUser(ctx, ip, sourceMes.Uid)
	if err != nil {
		err = fmt.Errorf("HandlePickIpMessage->SaveIpMapping err,%v", err)
		message.SendExceptionWebsocketMessage(ctx, conn)
		return
	}

	// 发送消息给客户端
	err = message.SendWebSocketMessage(ctx, conn, &message.ReplyMessageData{
		Data: &message.PickIpReplyMessage{
			Ip:  ip,
			Uid: sourceMes.Uid,
		},
	}, message.MesTypeReplyPickIp, "")
	if err != nil {
		log.E(ctx, err).Msgf("所有im节点均不可用")
		return
	}
	return
}

type BalanceRepo interface {
	GetUserIp(ctx context.Context, uid int) (ip string, err error)
	SaveIpUser(ctx context.Context, ip string, uid int) (err error)
	DeleteIpUser(ctx context.Context, uid int) (err error)
	DeleteIp(ctx context.Context, ip string) (err error)
}
