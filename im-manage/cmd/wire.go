//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"undersea/im-manage/conf"
	"undersea/im-manage/internal/biz"
	"undersea/im-manage/internal/data"
	server2 "undersea/im-manage/internal/server"
	"undersea/im-manage/internal/service"
)

func initApp() (*app, error) {
	panic(wire.Build(
		context.Background,
		conf.NewConf,
		newApp,
		server2.NewWebsocketServer,
		server2.NewGrpcClient,
		service.NewManageService,
		service.NewLoginService,
		biz.NewLoginUseCase,
		data.NewLoginRepo,
		data.NewRedisClient,
	))
}
