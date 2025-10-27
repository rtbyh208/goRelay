package relaytcpserver

import (
	"crypto/cipher"
	"fmt"
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"io"
	"net"
	"time"
)

func RunTcpServer(addr string, whitelist []string, id string, directConn string, registerID string, aead cipher.AEAD) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		goLog.Error("listen ", addr, " error ", err)
		return
	}
	goLog.Info("listen ", listen.Addr().String(), " successful")

	for {
		conn, err := listen.Accept()
		if err != nil {
			goLog.Error("accept error", err)
			continue
		}
		go worker(conn, whitelist, id, directConn, registerID, aead)
	}
}

func worker(conn net.Conn, whitelist []string, id string, directConn string, registerID string, aead cipher.AEAD) {
	goLog.Debug(conn.RemoteAddr().String(), " is connection")

	defer conn.Close()

	if !pkg.IsWhitelisted(conn.RemoteAddr().String(), whitelist) {
		goLog.Debug(conn.RemoteAddr().String(), " not while list")
		return
	}

	func() {
		clientMapLock.Lock()
		defer clientMapLock.Unlock()
		var cInfo pipeprotocol.ClientInfo
		cInfo.Conn = conn
		cInfo.Time = time.Now().Unix()
		clientMap[conn.RemoteAddr().String()] = cInfo
	}()
	defer func() {
		clientMapLock.Lock()
		defer clientMapLock.Unlock()
		delete(clientMap, conn.RemoteAddr().String())
	}()

	conn.SetReadDeadline(time.Now().Add(time.Hour * 12))
	conn.SetWriteDeadline(time.Now().Add(time.Hour * 12))

	for {
		buf := make([]byte, pipeprotocol.MaxPackageLen)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				goLog.Error("conn read relay server error,error: ", err)
			}

			break
		}

		msg := buf[:n]
		var p pipeprotocol.ClientProtocolInfo
		p.Id = id
		p.DirectConn = directConn
		p.Conn = conn.RemoteAddr().String()
		p.Counts = fmt.Sprintf("%d", time.Now().UnixNano()+2152384)

		goLog.Debug("encode onec keys: ", p.Counts)

		aesOnec := pipeprotocol.AesNewNonece(p.Counts)
		aesEnBuf := pipeprotocol.AesEncode(aead, aesOnec, msg)

		p.Buf = append(p.Buf, aesEnBuf...)
		p.RegisterID = registerID

		goLog.Debug("id: ", p.Id, " conn: ", conn.RemoteAddr().String(), " registerID: ", p.RegisterID)

		jsonBuf, err := pkg.JsonMarshal(p)

		if err != nil {
			goLog.Error("json marshal error", err)
			return
		}

		func() {
			pipeLock.Lock()
			defer pipeLock.Unlock()
			if nil == pipeClient {
				goLog.Error("pipeClientConn is nil")
				return
			}
			// goLog.Debug(p.Conn, " The message has been sent to the pipeline.")
			pipeprotocol.SendMessage(pipeClient, jsonBuf)
		}()

		cInfo := func(pConn string) pipeprotocol.ClientInfo {
			clientMapLock.Lock()
			defer clientMapLock.Unlock()

			return clientMap[pConn]
		}(conn.RemoteAddr().String())

		cInfo.Time = time.Now().Unix()

		func() {
			clientMapLock.Lock()
			defer clientMapLock.Unlock()
			clientMap[conn.RemoteAddr().String()] = cInfo
		}()
	}
}
