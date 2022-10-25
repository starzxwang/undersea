// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/biz"
	"undersea/im_balance/internal/server"
	"undersea/im_balance/internal/service"
)

// Injectors from wire.go:

func initApp() (*app, error) {
	contextContext := context.Background()
	confConf := conf.NewConf()
	imManagerUseCase := biz.NewImManagerUseCase(confConf)
	imManagerService := service.NewImManagerService(confConf, imManagerUseCase)
	tcpServer := server.NewTcpServer(confConf, imManagerService)
	balanceUseCase := biz.NewBalanceUseCase(confConf)
	balanceService := service.NewBalanceService(contextContext, balanceUseCase)
	websocketServer := server.NewWebsocketServer(confConf, balanceService)
	mainApp := newApp(contextContext, tcpServer, websocketServer)
	return mainApp, nil
}
