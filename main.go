package main

import (
	"fmt"
	pipeserver_main "goRelay/pipeServer"
	"goRelay/pkg"
	relayclient_main "goRelay/relayClient"
	relayserver_main "goRelay/relayServer"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "goRelay",
		Usage: "is a TCP-based intranet penetration tool written in go.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "show version",
			},
		},
		Commands: []*cli.Command{
			pipeserver_main.RunPipeServerCommand(),
			relayclient_main.RunRelayClientCommand(),
			relayserver_main.RunRelayServerCommand(),
		},
		Action: func(ctx *cli.Context) error {

			if ctx.Bool("version") {
				fmt.Println("version:", pkg.Version)
				fmt.Println("buildAt:", pkg.BuildAt)
				fmt.Println("gitCommit:", pkg.GitCommit)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
