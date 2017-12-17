package main

import (
	"bytes"
	"net/http"
	"os/exec"

	"github.com/kataras/iris"
)

func (s *server) registerEndpoints() {
	s.HTTPServer.Get("/config", s.getConfig)
	s.HTTPServer.Get("/cmd/{commandName}", s.execCommand)
}

func (s *server) getConfig(ctx iris.Context) {
	ctx.JSON(s.Config)
}

func (s *server) execCommand(ctx iris.Context) {
	s.Logger.Errorf("execCommand")
	commandName := ctx.Params().Get("commandName")

	cmdConf, err := s.Config.getCommand(commandName)
	if err != nil {
		s.Logger.Errorf("execCommand failed error: %+v", err)
		ctx.StatusCode(http.StatusNotFound)
		return
	}

	stdin := bytes.NewBufferString(cmdConf.Stdin)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command(cmdConf.Cmd[0], cmdConf.Cmd[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	s.GetStore().Set(s.getRequestID(ctx),
		commandStoreItem{
			command: &cmdConf,
			stdin:   stdin,
			stdout:  stdout,
			stderr:  stderr,
		})
	s.Logger.Infof("store size=%+v", s.GetStore().Size())

	if cmdConf.Long {
		go func() {
			err := cmd.Start()
			if err != nil {
				s.Logger.Errorf("async cmd run fail id=%+v cmd=%+v err=%+v ", s.getRequestID(ctx), cmdConf.Name, err)
			}
			s.Logger.Debugf("async cmd run success id=%+v cmd=%+v", s.getRequestID(ctx), cmdConf.Name)
		}()
		ctx.Text(s.getRequestID(ctx))
		return
	}
	err = cmd.Run()
	if err != nil {
		s.Logger.Errorf("cmd run fail id=%+v cmd=%+v err=%+v ", s.getRequestID(ctx), cmdConf.Name, err)
		ctx.Text(stderr.String())
	}
	s.Logger.Debugf("cmd run success id=%+v cmd=%+v", s.getRequestID(ctx), cmdConf.Name)
	ctx.Text(stdout.String())

}
