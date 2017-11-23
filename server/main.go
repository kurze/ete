package main

import (
	"bytes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"os/exec"
)

func main() {
	server := createServer()
	log := server.Logger()
	c, err := readConfig("conf.yml")
	if err != nil {
		log.Error(err)
	}
	log.Debug("readConfig", c)

	declareCommand(c, server)

	server.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func declareCommand(config *Config, server *iris.Application) {
	log := server.Logger()
	for i, c := range config.Commands {
		log.Debugf("%v %+v", i, c)
		execCommand(server, c)
	}
}

func execCommand(server *iris.Application, c Command) {
	server.Get("/" + c.Name, func(ctx iris.Context) {
		cmd := exec.Command(c.Cmd[0], c.Cmd[1:]...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			ctx.Text(err.Error())
		}
		ctx.Text(out.String())
	})
}

func createServer() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Logger().SetLevel("debug")

	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})
	return app
}
