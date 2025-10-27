package pipeserver

import (
	"goRelay/pkg"
	"net"
)

func init() {
	goLog = pkg.NewLogger()
	clientConnMap = make(map[string]net.Conn, 0)
}
