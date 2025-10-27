package relaytcpserver

import (
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"net"
	"sync"
)

var (
	pipeClient      net.Conn
	err             error
	clientMap       map[string]pipeprotocol.ClientInfo
	pipeLock        sync.Mutex
	clientMapLock   sync.Mutex
	goLog           *pkg.Logger
	sleepTimeSecond int = 3
)
