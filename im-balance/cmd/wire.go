//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"undersea/im-balance/conf"
	"undersea/im-balance/internal/biz"
	"undersea/im-balance/internal/data"
	server2 "undersea/im-balance/internal/server"
	"undersea/im-balance/internal/service"
)

func initApp() (*app, error) {
	panic(wire.Build(
		context.Background,
		conf.NewConf,
		newApp,
		server2.NewWebsocketServer,
		server2.NewGrpcServer,
		service.NewBalanceService,
		service.NewHeartBeatService,
		biz.NewHeartBeatUseCase,
		biz.NewBalanceUseCase,
		data.NewRedisClient,
		data.NewBalanceRepo,
	))
}
