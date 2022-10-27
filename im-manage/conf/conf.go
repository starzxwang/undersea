package conf

import (
	"time"
	"undersea/pkg/viper"
)

type Conf struct {
	Ws struct {
		Addr              string // websocket host
		HeartBeatInterval time.Duration
	}

	Grpc struct {
		ImBalanceAddr     string
		HeartBeatInterval time.Duration
	}
}

func NewConf() Conf {
	return Conf{
		Ws: struct {
			Addr              string
			HeartBeatInterval time.Duration
		}{
			Addr:              viper.V().GetString("im-manage.ws.addr"),
			HeartBeatInterval: time.Duration(viper.V().GetInt("im-manage.ws.heart_beat.interval")) * time.Second,
		},
		Grpc: struct {
			ImBalanceAddr     string
			HeartBeatInterval time.Duration
		}{
			ImBalanceAddr:     viper.V().GetString("im-balance.grpc.addr"),
			HeartBeatInterval: time.Duration(viper.V().GetInt("im-balance.grpc.heart_beat.interval")) * time.Second,
		},
	}
}
