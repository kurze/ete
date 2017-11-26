package main

import (
	"os/exec"

	"github.com/kataras/iris"
)

func (s *server) getConfig(ctx iris.Context) {
	ctx.JSON(s.Config)
}

func (s *server) registerCommand(c command) {
	s.HTTPServer.Get("/cmd/"+c.Name, func(ctx iris.Context) {

		out, err := exec.Command(c.Cmd[0], c.Cmd[1:]...).CombinedOutput()

		if err != nil {
			s.Logger.Error("cmd run fail ", err)

			if _, err := ctx.Text(err.Error()); err != nil {
				s.Logger.Error("send error fail ", err)
			}
		}

		if _, err := ctx.Text(string(out)); err != nil {
			s.Logger.Error("send result fail ", err)
		}
	})
}
