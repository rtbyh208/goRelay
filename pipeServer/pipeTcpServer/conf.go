package pipeserver

import (
	"goRelay/pkg"
	"net"
	"sync"
)

var relayConn map[string]net.Conn
var clientConnMap map[string]net.Conn
var clientConnMapLock sync.Mutex
var sendLock sync.Mutex
var goLog *pkg.Logger
