package biz

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"undersea/im-manage/conf"
	"undersea/pkg/message"
	"undersea/pkg/util"
)

var (
	userConnMapping sync.Map
)

type LoginUseCase struct {
	conf      conf.Conf
	loginRepo LoginRepo
}

type HeartBeatItem struct {
	pingChan chan struct{}
	conn     *websocket.Conn
}

func NewLoginUseCase(loginRepo LoginRepo) *LoginUseCase {
	return &LoginUseCase{
		loginRepo: loginRepo,
	}
}

func (uc *LoginUseCase) Login(ctx context.Context, mes *message.LoginMessage, conn *websocket.Conn) (err error) {
	// 这个用户是不是属于这一台机器
	userIp := uc.loginRepo.GetUserIp(ctx, mes.Uid)
	if userIp == "" {
		err = fmt.Errorf("该用户节点已经变更")
		return
	}

	if util.GetIpAddr() != userIp {
		err = fmt.Errorf("用户节点不匹配")
		return
	}

	// 这个用户是否已经登录
	if v, ok := userConnMapping.Load(mes.Uid); ok {
		// 用户已经登录，维持心跳
		v.(*HeartBeatItem).pingChan <- struct{}{}
		return
	}

	// 建立uid和conn的映射关系
	heartBeatItem := &HeartBeatItem{
		pingChan: make(chan struct{}, 1),
		conn:     conn,
	}
	userConnMapping.Store(mes.Uid, heartBeatItem)
	go uc.heartBeat(mes.Uid, heartBeatItem)
	heartBeatItem.pingChan <- struct{}{}

	return
}

func (uc *LoginUseCase) heartBeat(uid int, heartBeatItem *HeartBeatItem) {
	online := true
	for online {
		select {
		case <-heartBeatItem.pingChan:
			heartBeatItem.conn.SetReadDeadline(time.Now().Add(uc.conf.Ws.HeartBeatInterval))
		case <-time.After(uc.conf.Ws.HeartBeatInterval):
			heartBeatItem.conn.Close()
			userConnMapping.Delete(uid)
			online = false
		}
	}
}

type LoginRepo interface {
	GetUserIp(ctx context.Context, uid int) (ip string)
}
