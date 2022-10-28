//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"undersea/im-api/conf"
	"undersea/im-api/internal/biz"
	"undersea/im-api/internal/data"
	server2 "undersea/im-api/internal/server"
	"undersea/im-api/internal/service"
)

func initApp() (*app, error) {
	panic(wire.Build(
		context.Background,
		conf.NewConf,
		newApp,
		server2.NewHttpServer,
		service.NewUserService,
		biz.NewUserUseCase,
		data.NewMysql,
		data.NewUserRepo,
	))
}
