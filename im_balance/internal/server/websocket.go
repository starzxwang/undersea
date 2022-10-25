package server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/service"
	"undersea/pkg/log"
)

var (
	upgrader = websocket.Upgrader{
		//允许跨域访问
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebsocketServer struct {
	conf           conf.Conf
	balanceService *service.BalanceService
}

func NewWebsocketServer(conf conf.Conf, balanceService *service.BalanceService) *WebsocketServer {
	return &WebsocketServer{conf: conf, balanceService: balanceService}
}

func (s *WebsocketServer) Name() string {
	return "websocket"
}

func (s *WebsocketServer) Start(ctx context.Context) (err error) {
	// 监听来自前端的websocket
	http.HandleFunc("/ws", s.wsHandler)

	//服务端启动
	log.I(ctx).Msgf("负载均衡模块开始监听websocket端口：%s", s.conf.WsAddr)
	err = http.ListenAndServe(s.conf.WsAddr, nil)
	if err != nil {
		log.E(ctx, err).Msgf("websocket start err")
		return
	}

	return
}

func (s *WebsocketServer) wsHandler(w http.ResponseWriter, r *http.Request) {
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
		go s.balanceService.HandleClientMessage(conn, data)
	}
}

func (s *WebsocketServer) Stop(ctx context.Context) error {
	return nil
}
