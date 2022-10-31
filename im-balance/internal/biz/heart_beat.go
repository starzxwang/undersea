package biz

import (
	"context"
	"sync"
	"time"
	"undersea/im-balance/conf"
	"undersea/pkg/log"
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
	conf        conf.Conf
	balanceRepo BalanceRepo
}

func NewHeartBeatUseCase(conf conf.Conf, balanceRepo BalanceRepo) (uc *HeartBeatUseCase) {
	return &HeartBeatUseCase{conf: conf, balanceRepo: balanceRepo}
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
		go uc.listenImHeartBeat(mes, heartBeatItem)
		return
	}

	v2, ok := v.(*sync.Map).Load(mes.Ip)
	if ok {
		// 心跳
		v2.(*HeartBeatItem).pingChan <- struct{}{}
	} else {
		// 开启协程监听这个im节点
		heartBeatItem := &HeartBeatItem{
			pingChan:     make(chan struct{}, 1),
			lastPingTime: time.Now(),
		}
		v.(*sync.Map).Store(mes.Ip, heartBeatItem)
		imServer.ipMap.Store(mes.ServiceName, v)
		go uc.listenImHeartBeat(mes, heartBeatItem)
	}

	return
}

// 监听当前im节点的心跳
func (uc *HeartBeatUseCase) listenImHeartBeat(mes *message.HeartBeatMessage, heartBeatItem *HeartBeatItem) {
	ctx := context.Background()
	online := true
	for online {
		select {
		case <-heartBeatItem.pingChan:
			heartBeatItem.lastPingTime = time.Now()
		default:
			if time.Now().Sub(heartBeatItem.lastPingTime) > uc.conf.Grpc.HeartBeatInterval {
				// 删除redis和本地缓存
				err := uc.balanceRepo.DeleteIp(ctx, mes.Ip)
				if err != nil {
					log.E(ctx, err).Msgf("listenImHeartBeat->DeleteIp err")
					time.Sleep(3 * time.Second)
					continue
				}
				imServer.ipMap.Delete(mes.ServiceName)
				online = false
			}
		}

		time.Sleep(time.Millisecond)
	}
}
