package biz

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"time"
	"undersea/im_balance/conf"
	"undersea/pkg/log"
	"undersea/pkg/message"
)

var (
	imServer = &imManagerServer{}
)

type imManagerServer struct {
	ipList       map[string]chan struct{} // 当前可用的im节点列表
	currentIndex int                      // 当前请求所选ip的索引值
	rwLock       sync.RWMutex
}

type ImManagerUseCase struct {
	conf conf.Conf
}

func NewImManagerUseCase(conf conf.Conf) (uc *ImManagerUseCase) {
	return &ImManagerUseCase{conf: conf}
}

func (uc *ImManagerUseCase) HandleHeartBeatMessage(ctx context.Context, conn net.Conn, mes *message.Message) (err error) {
	var sourceMes *message.HeartBeatMessage
	err = json.Unmarshal([]byte(mes.Data), &sourceMes)
	if err != nil {
		log.E(ctx, err).Msgf("json.Unmarshal err")
		return
	}

	// 注册或更新服务节点
	uc.SaveNode(sourceMes.ImManagerIp, conn)
	return
}

// 注册或更新服务节点
func (uc *ImManagerUseCase) SaveNode(ip string, conn net.Conn) {
	imServer.rwLock.RLock()
	v, ok := imServer.ipList[ip]
	imServer.rwLock.RUnlock()
	if ok {
		// 心跳
		v <- struct{}{}
	} else {
		// 开启协程监听这个im节点
		imServer.rwLock.Lock()
		imServer.ipList[ip] = make(chan struct{}, 1)
		imServer.rwLock.Unlock()
		go uc.listenImHeartBeat(ip, conn)
	}

	return
}

// 监听当前im节点的心跳
func (uc *ImManagerUseCase) listenImHeartBeat(ip string, conn net.Conn) {
	ctx := context.Background()
	online := true
	for online {
		imServer.rwLock.RLock()
		nodeChan := imServer.ipList[ip]
		imServer.rwLock.RUnlock()
		select {
		case <-nodeChan:
			err := conn.SetDeadline(time.Now().Add(uc.conf.ConnActiveTime))
			if err != nil {
				log.E(ctx, err).Msgf("listenImNodeChan->conn.SetDeadline err")
				continue
			}
		case <-time.After(uc.conf.ConnActiveTime):
			imServer.rwLock.Lock()
			delete(imServer.ipList, ip)
			imServer.rwLock.Unlock()

			// 关闭连接
			conn.Close()

			// 退出监听协程
			online = false
		}
	}
}
