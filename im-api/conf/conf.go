package conf

import (
	"time"
	"undersea/pkg/viper"
)

type Conf struct {
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

	MySQL struct {
		Username     string
		Password     string
		DbName       string
		Host         string
		Port         int
		Charset      string
		Zone         string
		MaxOpenConns int
		MaxIdleConns int
	}

	Http struct {
		Addr string
	}
}

func NewConf() Conf {
	return Conf{
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
		MySQL: struct {
			Username     string
			Password     string
			DbName       string
			Host         string
			Port         int
			Charset      string
			Zone         string
			MaxOpenConns int
			MaxIdleConns int
		}{
			Username:     viper.V().GetString("mysql.username"),
			Password:     viper.V().GetString("mysql.password"),
			DbName:       viper.V().GetString("mysql.db_name"),
			Host:         viper.V().GetString("mysql.host"),
			Port:         viper.V().GetInt("mysql.port"),
			Charset:      viper.V().GetString("mysql.charset"),
			Zone:         viper.V().GetString("mysql.zone"),
			MaxOpenConns: viper.V().GetInt("mysql.max_open_conns"),
			MaxIdleConns: viper.V().GetInt("mysql.max_idle_conns"),
		},
		Http: struct {
			Addr string
		}{
			Addr: viper.V().GetString("im_api.http.addr"),
		},
	}
}
