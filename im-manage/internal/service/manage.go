package service

import (
	"context"
	"github.com/gorilla/websocket"
	"undersea/pkg/log"
	"undersea/pkg/message"
)

type ManageService struct {
	loginService *LoginService
	//peerService  *PeerService
	//groupService *GroupService
}

func NewManageService(loginService *LoginService) *ManageService {
	return &ManageService{
		loginService: loginService,
		//peerService:  peerService,
		//groupService: groupService,
	}
}

func (s *ManageService) HandleClientMessage(conn *websocket.Conn, data []byte) {
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
	case message.MesTypeLogin:
		// 客户端需要一个im节点ip
		err = s.loginService.HandleLoginMessage(ctx, mes, conn)
	default:
	}

	if err != nil {
		log.E(ctx, err).Msgf("handle client msg err,%s", mes.Type)
		message.SendExceptionWebsocketMessage(ctx, conn)
		return
	}

	return
}
