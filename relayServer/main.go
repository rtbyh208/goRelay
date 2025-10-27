package relayserver_main

import (
	"errors"
	"fmt"
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	relaytcpserver "goRelay/relayServer/relayTcpServer"

	"github.com/urfave/cli/v2"
)

func RunRelayServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "relayServer",
		Usage: "relay server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf",
				Usage: "config path",
			},
			&cli.StringFlag{
				Name:  "directConn",
				Usage: "direct conn",
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

			if ctx.String("directConn") != "" {
				config.DirectConn = ctx.String("directConn")
			}

			pipeprotocol.Keys = append(pipeprotocol.Keys, pkg.IDHash(config.Key))

			RunRelayServer(config)

			return nil
		},
	}
}

func RunRelayServer(config Config) {

	goLog := pkg.NewLogger()
	if config.DebugLog {
		goLog.SetLogger(pkg.DebugLevel)
	} else {
		goLog.SetLogger(pkg.LogLevel)
	}

	configID := pkg.IDHash(config.Id)

	goLog.Debug("aes new cipher ids: ", configID)
	aead, err := pipeprotocol.AesNewCipher(configID)
	if err != nil {
		panic(err)
	}

	goLog.Debug("config: ", config)
	go relaytcpserver.ConnectPipeServer(config.PipeServerAddr, config.RegisterID, aead)

	relaytcpserver.RunTcpServer(config.ListenRelayServerAddr, config.WhiteIpList, configID, config.DirectConn, config.RegisterID, aead)
}
