package relayclient_main

import (
	"errors"
	"fmt"
	dbcache "goRelay/dbCache"
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	relaytcpclient "goRelay/relayClient/relayTcpClient"

	"github.com/urfave/cli/v2"
)

func RunRelayClientCommand() *cli.Command {
	return &cli.Command{
		Name:  "relayClient",
		Usage: "relay client",
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

			RunRelayClientServer(config)

			return nil
		},
	}
}

func RunRelayClientServer(config Config) {

	goLog := pkg.NewLogger()
	if config.DebugLog {
		goLog.SetLogger(pkg.DebugLevel)
	} else {
		goLog.SetLogger(pkg.LogLevel)
	}

	cache := dbcache.Init("cache.db")

	go cache.AutoSave()

	for _, v := range config.RealServerInfo {
		configID := pkg.IDHash(v.ID)
		cache.Set(configID, v.RealServerAddr)
	}

	goLog.Debug("config: ", config)

	if config.HttpServer != "" {
		go relaytcpclient.StartHttpServer(config.Key, config.HttpServer, config.PipeServerAddr, config.RegisterID, cache)
	}
	relaytcpclient.ConnectToPipeServer(config.PipeServerAddr, config.RegisterID, cache)
}
