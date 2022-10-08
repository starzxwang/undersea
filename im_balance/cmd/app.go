package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/service"
	"undersea/pkg/log"
)

type app struct {
	ctx              context.Context
	conf             conf.Conf
	balanceService   *service.BalanceService
	imManagerService *service.ImManagerService
}

var (
	upgrader = websocket.Upgrader{
		//允许跨域访问
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func newApp(ctx context.Context, conf conf.Conf, balanceService *service.BalanceService, imManagerService *service.ImManagerService) *app {
	return &app{
		ctx:              ctx,
		conf:             conf,
		balanceService:   balanceService,
		imManagerService: imManagerService,
	}
}

func (app *app) run() {
	// 监听来自im的服务注册和心跳
	go app.imManagerService.ListenTCP()

	// 监听来自前端的websocket
	http.HandleFunc("/ws", app.wsHandler)

	//服务端启动
	log.I(app.ctx).Msgf("负载均衡模块开始监听websocket端口：%s", app.conf.WsAddr)
	http.ListenAndServe(app.conf.WsAddr, nil)
}

// 监听Im_manager可用节点
func (app *app) listenIMNodes() {
	for {
		//
		time.Sleep(time.Millisecond)
	}
}

func (app *app) wsHandler(w http.ResponseWriter, r *http.Request) {
	//收到http请求(upgrade),完成websocket协议转换
	//在应答的header中放上upgrade:websoket
	var (
		conn *websocket.Conn
		err  error
		//msgType int
		data []byte
		ctx  context.Context
	)
	if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		//报错了，直接返回底层的websocket链接就会终断掉
		fmt.Println("wsHandler:get conn err=", err)
		return
	}

	//设置连接超时时间
	//conn.SetReadDeadline(time.Now().Add(time.Duration(common.Config.Heart_beat_timeout) * time.Second))

	defer conn.Close()

	//得到了websocket.Conn长连接的对象，实现数据的收发
	for {
		log.I(ctx).Msgf("wsHandler:等待客户端连接")

		if _, data, err = conn.ReadMessage(); err != nil {
			//报错关闭websocket
			log.I(ctx).Msgf("wsHandler:conn.ReadMessage() err=%v", err)
			return
		}

		log.I(ctx).Msgf("wsHandler:接收到客户端消息，msg=%s", string(data))

		//开启协程，处理接收到的消息
		go app.balanceService.HandleClientMessage(conn, data)
	}
}
