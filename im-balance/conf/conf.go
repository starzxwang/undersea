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

	Redis struct {
		Addr      string
		UserName  string
		Password  string
		BalanceDb int
	}
}

func NewConf() Conf {
	return Conf{
		Ws: struct {
			Addr string
		}{
			Addr: viper.V().GetString("im-balance.ws_addr"),
		},
		Grpc: struct {
			Addr              string
			HeartBeatInterval time.Duration
		}{
			Addr:              viper.V().GetString("im-balance.grpc.addr"),
			HeartBeatInterval: time.Duration(viper.V().GetInt("im-balance.grpc.heart_beat.interval")) * time.Second,
		},
		Redis: struct {
			Addr      string
			UserName  string
			Password  string
			BalanceDb int
		}{
			Addr:      viper.V().GetString("redis.addr"),
			Password:  viper.V().GetString("redis.password"),
			BalanceDb: viper.V().GetInt("redis.balance_db"),
			UserName:  viper.V().GetString("redis.username"),
		},
	}
}
