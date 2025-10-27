package pipeserver

import (
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"net"
)

func worker(conn net.Conn, whiteIpList []string, blackIpList []string) {

	goLog.Info("Client ", conn.RemoteAddr().String(), " is connected.")

	defer func() {
		goLog.Info("Client ", conn.RemoteAddr().String(), " closed")
		conn.Close()
	}()

	if relayConn == nil {
		relayConn = make(map[string]net.Conn)
	}

	if pkg.IsBlacklisted(conn.RemoteAddr().String(), blackIpList) {
		goLog.Debug(conn.RemoteAddr().String(), " is black list")
		return
	}

	for {
		msg, err := pipeprotocol.RecvMessgae(conn)
		if err != nil {
			goLog.Error("from ", conn.RemoteAddr().String(), " recv messages error: ", err)
			break
		}

		var p pipeprotocol.ClientProtocolInfo

		err = pkg.JsonUnmarshal(msg, &p)
		if err != nil {
			goLog.Error("recv message error", err)
			return
		}

		if string(p.Buf) == "relayConn" {
			relayWorker(conn, p.RegisterID)
			continue
		}

		if nil == relayConn[p.RegisterID] {
			goLog.Error("relayConn id: ", p.RegisterID, " relayConn is nil", " relayConnPools: ", relayConn)
			break
		}
		goLog.Debug("received ", conn.RemoteAddr().String(), " is client worker.", " registerID:", p.RegisterID, " data len: ", len(msg))

		if err != nil {
			goLog.Error("Error: JSON unmarshal failed:", err, " jsonData:", msg)
			break
		}

		func() {
			clientConnMapLock.Lock()
			defer clientConnMapLock.Unlock()

			clientConnMap[p.Id] = conn
			// goLog.Debug(p.Conn, " has received a client request\nBuf:\n", string(p.Buf), "\nconn:\n", p.Conn, "\nid:", p.Id)
			// goLog.Debug("clientConnMaps: ", clientConnMap)
		}()

		func() {
			sendLock.Lock()
			defer sendLock.Unlock()
			pipeprotocol.SendMessage(relayConn[p.RegisterID], msg)
		}()
	}
}
