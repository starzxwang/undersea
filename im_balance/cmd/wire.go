// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"undersea/im_balance/conf"
	"undersea/im_balance/internal/biz"
	"undersea/im_balance/internal/service"
)

func initApp() (*app, error) {
	panic(wire.Build(
		context.Background,
		conf.NewConf,
		newApp,
		service.NewBalanceService,
		service.NewImManagerService,
		biz.NewImManagerUseCase,
		biz.NewBalanceUseCase,
	))
}
