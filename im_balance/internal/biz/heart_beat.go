package biz

import (
	"sync"
	"time"
	"undersea/im_balance/conf"
	"undersea/pkg/message"
)

var (
	imServer = &imManagerServer{}
)

type imManagerServer struct {
	ipMap sync.Map // 当前可用的im节点列表
}

type HeartBeatItem struct {
	pingChan     chan struct{}
	lastPingTime time.Time
}

type HeartBeatUseCase struct {
	conf conf.Conf
}

func NewHeartBeatUseCase(conf conf.Conf) (uc *HeartBeatUseCase) {
	return &HeartBeatUseCase{conf: conf}
}

// 注册或更新服务节点
func (uc *HeartBeatUseCase) SaveHeartBeat(mes *message.HeartBeatMessage) {
	v, ok := imServer.ipMap.Load(mes.ServiceName)
	if !ok {
		var serviceIpMap sync.Map
		heartBeatItem := &HeartBeatItem{
			pingChan:     make(chan struct{}, 1),
			lastPingTime: time.Now(),
		}
		serviceIpMap.Store(mes.Ip, heartBeatItem)
		imServer.ipMap.Store(mes.ServiceName, &serviceIpMap)
		go uc.listenImHeartBeat(mes.Ip, heartBeatItem)
		return
	}

	v, ok = v.(*sync.Map).Load(mes.Ip)
	if ok {
		// 心跳
		v.(chan struct{}) <- struct{}{}
	} else {
		// 开启协程监听这个im节点
		heartBeatItem := &HeartBeatItem{
			pingChan:     make(chan struct{}, 1),
			lastPingTime: time.Now(),
		}
		imServer.ipMap.Store(mes.Ip, heartBeatItem)
		go uc.listenImHeartBeat(mes.Ip, heartBeatItem)
	}

	return
}

// 监听当前im节点的心跳
func (uc *HeartBeatUseCase) listenImHeartBeat(ip string, heartBeatItem *HeartBeatItem) {
	online := true
	for online {
		select {
		case <-heartBeatItem.pingChan:
			heartBeatItem.lastPingTime = time.Now()
		default:
			if time.Now().Sub(heartBeatItem.lastPingTime) > uc.conf.Grpc.HeartBeatInterval {
				imServer.ipMap.Delete(ip)
				online = false
			}
		}

		time.Sleep(time.Millisecond)
	}
}
