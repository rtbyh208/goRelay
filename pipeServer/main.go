package pipeserver_main

import (
	"errors"
	"fmt"
	pipeprotocol "goRelay/pipeProtocol"
	pipeserver "goRelay/pipeServer/pipeTcpServer"
	"goRelay/pkg"

	"github.com/urfave/cli/v2"
)

func RunPipeServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "pipeServer",
		Usage: "pipe server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf",
				Usage: "config path",
			},
		},
		Action: func(ctx *cli.Context) error {

			var config Config
			if ctx.String("conf") != "" {
				if pkg.LoadConfig(ctx.String("conf"), &config) != nil {
					fmt.Println("read config file error")
					return errors.New("read config file error")
				}
			}

			pipeprotocol.Keys = append(pipeprotocol.Keys, pkg.IDHash(config.Key))

			RunPipeServer(config)
			return nil
		},
	}
}

func RunPipeServer(config Config) {

	goLog := pkg.NewLogger()
	if config.DebugLog {
		goLog.SetLogger(pkg.DebugLevel)
	} else {
		goLog.SetLogger(pkg.LogLevel)
	}

	goLog.Debug("config: ", config)
	pipeserver.ListenTcpServer(config.ListenPipeServerAddr, config.WhiteIpList, config.BlackIpList)
}
