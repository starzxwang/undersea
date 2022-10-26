//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/biz"
	server2 "undersea/im_balance/internal/server"
	"undersea/im_balance/internal/service"
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
	))
}
