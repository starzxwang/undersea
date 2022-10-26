package conf

import (
	"time"
	"undersea/pkg/viper"
)

type Conf struct {
	Ws struct {
		Addr string // websocket host
	}

	Grpc struct {
		Addr              string
		HeartBeatInterval time.Duration
	}
}

func NewConf() Conf {
	return Conf{
		Ws: struct {
			Addr string
		}{
			Addr: viper.V().GetString("im_balance.ws_addr"),
		},
		Grpc: struct {
			Addr              string
			HeartBeatInterval time.Duration
		}{
			Addr:              viper.V().GetString("grpc.addr"),
			HeartBeatInterval: time.Duration(viper.V().GetInt("grpc.heart_beat.interval")) * time.Second,
		},
	}
}
