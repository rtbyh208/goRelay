package relaytcpclient

import (
	"crypto/cipher"
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"net"
	"sync"
)

var (
	clientConnections map[string]pipeprotocol.ClientInfo
	clientAEAD        map[string]cipher.AEAD
	goLog             *pkg.Logger
	clientMapMutex    sync.Mutex
	relayConn         net.Conn
	pipeClientConn    net.Conn
	err               error
	sleepTimeSec      int = 3
)
