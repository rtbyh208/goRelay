package relaytcpclient

import (
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"net"
	"time"
)

func RunTcpServer(addr string, whitelist []string, realServerInfoMap map[string]string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		goLog.Error("listen error: ", err)
		return
	}
	goLog.Info("listen ", addr, " successful")

	for {
		relayConn, err = listen.Accept()
		if err != nil {
			goLog.Error("accept error", err)
		}

		worker(relayConn, whitelist, realServerInfoMap)
	}
}

func worker(conn net.Conn, whitelist []string, realServerInfoMap map[string]string) {
	goLog.Info("client ", conn.RemoteAddr().String(), " is connected")

	defer conn.Close()

	if !pkg.IsWhitelisted(conn.RemoteAddr().String(), whitelist) {
		goLog.Debug(conn.RemoteAddr().String(), " not while list")
		return
	}

	for {
		msg, err := pipeprotocol.RecvMessgae(conn)
		if err != nil {
			goLog.Error("recv relay server error,error: ", err)
			break
		}

		var p pipeprotocol.ClientProtocolInfo
		err = pkg.JsonUnmarshal(msg, &p)
		if err != nil {
			goLog.Error("json unmarshal error ", err)
			continue
		}
		// goLog.Debug("from client ", conn.RemoteAddr().String(), " data , id: ", p.Id, " data len ", len(msg))

		_, isok := clientConnections[p.Conn]
		if !isok {
			realServerAddr, isok := realServerInfoMap[p.Id]
			if !isok {
				goLog.Error("not fount realservers ,id: ", p.Id)
				break
			}

			if realServerAddr == "directConn" {
				if p.CommandID != 100 {
					if p.DirectConn == "" {
						goLog.Error(p.Id, "it is a dynamic link, but no value has been set.")
						continue
					}
				}
				realServerAddr = p.DirectConn
			}
			if realServerAddr == "" {
				continue
			}

			realServer, err := net.Dial("tcp", realServerAddr)
			if err != nil {
				goLog.Error("dial real server error", err)
				break
			}

			var c pipeprotocol.ClientInfo
			c.Conn = realServer
			c.Time = time.Now().Unix()

			func() {
				clientMapMutex.Lock()
				defer clientMapMutex.Unlock()
				clientConnections[p.Conn] = c
			}()

			// go connectionRealServer(realServer, p.Conn, p.Id)

		}
		cInfo, _ := clientConnections[p.Conn]

		if p.CommandID == 100 {
			goLog.Info("close ", p.Id, " connecttion, commandid: ", p.CommandID)
			cInfo.Conn.Close()
			break
		}

		sendLen := 0
		for sendLen < len(p.Buf) {
			n, err := cInfo.Conn.Write(p.Buf[sendLen:])
			if err != nil {
				goLog.Error("to real server write data error: ", err)
				break
			}
			sendLen += n
		}

		func() {
			clientMapMutex.Lock()
			defer clientMapMutex.Unlock()
			clientConnections[p.Conn] = cInfo
		}()
	}
}
