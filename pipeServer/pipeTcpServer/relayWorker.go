package pipeserver

import (
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"net"
)

func relayWorker(conn net.Conn, registerID string) {
	// relayConn = conn
	relayConn[registerID] = conn

	var p pipeprotocol.ClientProtocolInfo
	p.Buf = []byte("isok")

	relayConnBuf, err := pkg.JsonMarshal(p)
	if err != nil {
		goLog.Error("relay worker error", err)
		return
	}
	pipeprotocol.SendMessage(conn, relayConnBuf)
	// defer func() {
	// 	// relayConn[registerID] = nil
	// 	delete(relayConn, registerID)
	// }()

	goLog.Info(conn.RemoteAddr().String(), " is relay worker.", "registerID:", registerID, " relayConn pool: ", relayConn)

	for {
		msg, err := pipeprotocol.RecvMessgae(conn)
		if err != nil {
			goLog.Error(conn.RemoteAddr(), " RecvMessgae error,err: ", err)
			break
		}

		var p pipeprotocol.ClientProtocolInfo

		func() {
			clientConnMapLock.Lock()
			defer clientConnMapLock.Unlock()
			clientConnMap[p.Id] = conn
		}()

		err = pkg.JsonUnmarshal(msg, &p)
		if err != nil {
			goLog.Error("Error: JSON unmarshal failed:", err, " jsonData:", msg)
			break
		}

		relayServerConn, isok := func() (net.Conn, bool) {
			clientConnMapLock.Lock()
			defer clientConnMapLock.Unlock()
			r, e := clientConnMap[p.Id]
			return r, e
		}()

		if !isok {
			goLog.Error("Error: not fount relayServer Conn")
			break
		}
		goLog.Debug("id: ", p.Id, " from: ", conn.RemoteAddr().String(), " messages", " data len: ", len(msg))

		pipeprotocol.SendMessage(relayServerConn, msg)
	}
}
