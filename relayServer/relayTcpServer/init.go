package relaytcpserver

import (
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
)

func init() {
	clientMap = make(map[string]pipeprotocol.ClientInfo, 0)
	goLog = pkg.NewLogger()
}
