package relaytcpclient

import (
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
)

func init() {
	clientConnections = make(map[string]pipeprotocol.ClientInfo, 0)
	goLog = pkg.NewLogger()
}
