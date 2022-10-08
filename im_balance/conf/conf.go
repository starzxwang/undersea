package conf

import (
	"time"
	"undersea/pkg/viper"
)

type Conf struct {
	WsAddr         string        // websocket host
	TcpAddr        string        // tcp host
	ConnActiveTime time.Duration // tcp连接最大空闲时间，超过了就会被关闭
}

func NewConf() Conf {
	return Conf{
		WsAddr:  viper.V().GetString("im_balance.ws_addr"),
		TcpAddr: viper.V().GetString("im_balance.tcp_addr"),
		ConnActiveTime: time.Duration(viper.V().GetInt("im_balance.conn_active_time")) * time.Second,
	}
}
