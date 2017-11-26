package main

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

type server struct {
	HTTPServer *iris.Application
	Logger     *golog.Logger
	Config     *config
}

func newServer(config *config) *server {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Logger().SetLevel(config.LogLevel)

	return &server{
		Logger:     app.Logger(),
		HTTPServer: app,
		Config:     config,
	}
}

func (s *server) run() {
	err := s.HTTPServer.Run(iris.Addr(s.Config.getListenAddress()), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		s.Logger.Fatal("http server fail", err)
	}
}

func (s *server) registerEndpoints() {
	s.HTTPServer.Get("config", s.getConfig)

	for i, c := range s.Config.Commands {
		s.Logger.Debugf("%v %+v", i, c)
		s.registerCommand(c)
	}
}
