package main

import (
	"encoding/base64"
	"math/rand"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

type server struct {
	HTTPServer *iris.Application
	Logger     *golog.Logger
	Config     *config
	store      *commandStore
	//ws         *websocket.Server
}

func (s *server) GetStore() *commandStore {
	return s.store
}

const keyRequestID = "reqID"

func newServer(config *config) *server {

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Status:            true,
		IP:                true,
		Method:            true,
		Path:              true,
		MessageContextKey: keyRequestID,
	}))
	app.Logger().SetLevel(config.LogLevel)

	s := server{
		Logger:     app.Logger(),
		HTTPServer: app,
		Config:     config,
		store:      newCommandStore(),
		//ws:         ws,
	}

	app.Use(s.before)

	//ws := websocket.New(websocket.Config{
	//	ReadBufferSize:  1024,
	//	WriteBufferSize: 1024,
	//})

	return &s
}

func (s *server) run() {
	err := s.HTTPServer.Run(iris.Addr(s.Config.getListenAddress()), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		s.Logger.Fatal("http server fail", err)
	}
}

func (s *server) generateRequestID() string {
	r := make([]byte, 6)
	_, err := rand.Read(r)
	if err != nil {
		s.Logger.Fatal("According to the lib code, *this* cannot be display.")
	}
	return base64.StdEncoding.EncodeToString(r)
}

func (s *server) before(ctx iris.Context) {
	ctx.Values().Set(keyRequestID, s.generateRequestID())
	ctx.Next()
}

func (s *server) getRequestID(ctx iris.Context) string {
	result, ok := ctx.Values().Get(keyRequestID).(string)
	if !ok {
		s.Logger.Fatal("someone do something wrong with requestID")
	}
	return result
}
