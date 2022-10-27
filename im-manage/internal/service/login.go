package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"undersea/im-manage/internal/biz"
	"undersea/pkg/message"
)

type LoginService struct {
	loginUseCase *biz.LoginUseCase
}

func NewLoginService(loginUseCase *biz.LoginUseCase) *LoginService {
	return &LoginService{
		loginUseCase: loginUseCase,
	}
}

// 处理登录消息
func (s *LoginService) HandleLoginMessage(ctx context.Context, mes *message.Message, conn *websocket.Conn) (err error) {
	// 解析消息
	var sourceMes *message.LoginMessage
	err = json.Unmarshal([]byte(mes.Data), &sourceMes)
	if err != nil {
		err = fmt.Errorf("HandleLoginMessage->json.Unmarshal err,%v", err)
		return
	}

	err = s.loginUseCase.Login(ctx, sourceMes, conn)
	if err != nil {
		err = fmt.Errorf("HandleLoginMessage->Login err")
		return
	}

	return
}
