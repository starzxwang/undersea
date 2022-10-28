package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	server2 "undersea/im-api/internal/server"
	"undersea/pkg/log"
)

type app struct {
	ctx     context.Context
	cancel  context.CancelFunc
	servers []server2.Server
}

func newApp(ctx context.Context, httpServer *server2.HttpServer) *app {
	ctx, cancel := context.WithCancel(context.Background())
	app := &app{
		ctx:    ctx,
		cancel: cancel,
	}
	app.servers = append(app.servers, httpServer)
	return app
}

func (app *app) run() (err error) {
	var wg sync.WaitGroup
	for _, s := range app.servers {
		wg.Add(1)
		go func(s server2.Server) {
			defer wg.Done()

			err = s.Start(app.ctx)
			if err != nil {
				log.E(app.ctx, err).Msgf("%s start failed", s.Name())
			}

			return
		}(s)

		wg.Add(1)
		go func(s server2.Server) {
			<-app.ctx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()
			err = s.Stop(ctx)
			if err != nil {
				log.E(ctx, err).Msgf("%s stop failed", s.Name())
			}

			wg.Done()
		}(s)
	}

	wg.Add(1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() error {
		defer wg.Done()
		select {
		case <-app.ctx.Done():
			return app.ctx.Err()
		case <-quit:
			app.cancel()
			return nil
		}
	}()

	wg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return
	}

	return
}
