package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"undersea/im-api/conf"
	"undersea/im-api/internal/service"
	"undersea/pkg/api"
	"undersea/pkg/log"
)

type HttpServer struct {
	conf        conf.Conf
	userService *service.UserService
}

func NewHttpServer(conf conf.Conf, userService *service.UserService) *HttpServer {
	return &HttpServer{
		conf:        conf,
		userService: userService,
	}
}

func (s *HttpServer) Name() string {
	return "api"
}

func (s *HttpServer) Start(ctx context.Context) error {
	e := gin.Default()
	baseGroup := e.Group("")
	baseGroup.Use(api.Cors())
	baseGroup.POST("/v1/login", s.userService.Login)
	baseGroup.POST("/v1/register", s.userService.Register)

	err := e.Run(s.conf.Http.Addr)
	if err != nil {
		log.E(ctx, err).Msgf("Start->api run err")
		return err
	}
	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	log.I(ctx).Msgf("[%s] server stopping", s.Name())
	return nil
}
