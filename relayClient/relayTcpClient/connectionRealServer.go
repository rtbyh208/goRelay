package relaytcpclient

import (
	"fmt"
	dbcache "goRelay/dbCache"
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"io"
	"net"
	"time"
)

func connectionRealServer(realServerConn net.Conn, pConn string, id string, cache *dbcache.Caches) {
	goLog.Debug("connection real server ", realServerConn.RemoteAddr().String())

	realServerConn.SetReadDeadline(time.Now().Add(time.Hour * 12))
	realServerConn.SetWriteDeadline(time.Now().Add(time.Hour * 12))

	defer func() {
		clientMapMutex.Lock()
		defer clientMapMutex.Unlock()
		delete(clientConnections, pConn)
	}()

	for {
		stopClientID := fmt.Sprintf("Stop_%s", id)
		if cache.Exists(stopClientID) {
			goLog.Debug("realservers is disabled ,id: ", id)
			return
		}

		buf := make([]byte, pipeprotocol.MaxPackageLen)

		n, err := realServerConn.Read(buf)
		if err != nil {
			if err != io.EOF {
				goLog.Error("from real server get data error", err)
			}
			func() {
				clientMapMutex.Lock()
				defer clientMapMutex.Unlock()
				goLog.Debug("delete client  ", pConn)
				delete(clientConnections, pConn)
			}()
			break
		}
		msg := buf[:n]

		var p pipeprotocol.ClientProtocolInfo
		p.Id = id
		p.Counts = fmt.Sprintf("%d", time.Now().UnixNano()-17601761)

		goLog.Debug("aes encode ", p.Counts)
		aesOnec := pipeprotocol.AesNewNonece(p.Counts)
		aesEnBuf := pipeprotocol.AesEncode(clientAEAD[p.Id], aesOnec, msg)

		p.Buf = append(p.Buf, aesEnBuf...)
		p.Conn = pConn

		jsonBuf, err := pkg.JsonMarshal(p)
		if err != nil {
			goLog.Error("json marshal error ", err)

			goLog.Info("delete client map ", pConn)
			delete(clientConnections, pConn)
			break
		}

		func() {
			clientMapMutex.Lock()
			defer clientMapMutex.Unlock()

			goLog.Debug("from ", realServerConn.RemoteAddr().String(), " data len: ", len(jsonBuf))
			pipeprotocol.SendMessage(relayConn, jsonBuf)
		}()

		AddRX(id, len(jsonBuf), cache)
	}
}
