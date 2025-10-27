package relaytcpserver

import (
	"crypto/cipher"
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"net"
	"time"
)

func ConnectPipeServer(pipeServerAddr string, registerID string, aead cipher.AEAD) {
	for {
		pipeClient, err = net.Dial("tcp", pipeServerAddr)
		if err != nil {
			goLog.Error("Unable to connect to the PipeServer server. ", err)
			pipeClient = nil
			time.Sleep(time.Duration(sleepTimeSecond) * time.Second)
			continue
		}
		goLog.Info("dial ", pipeServerAddr, " successful")

		relayWorker(pipeClient, registerID, aead)
	}
}

func relayWorker(pipeClient net.Conn, registerID string, aead cipher.AEAD) {
	for {
		msg, err := pipeprotocol.RecvMessgae(pipeClient)
		if err != nil {
			goLog.Error("recv mesage error,stop relay tcp connection. error: ", err)
			break
		}

		var p pipeprotocol.ClientProtocolInfo
		pkg.JsonUnmarshal(msg, &p)
		// goLog.Debug(p.Conn, " The message has returned.")
		// goLog.Debug("p.id: ", p.Id)

		cInfo, isok := func(pConn string) (pipeprotocol.ClientInfo, bool) {
			clientMapLock.Lock()
			defer clientMapLock.Unlock()
			cINfo, isok := clientMap[pConn]
			return cINfo, isok
		}(p.Conn)

		if !isok {
			goLog.Debug("not fount client conn,error,conn: ", p.Conn)

			var recallP pipeprotocol.ClientProtocolInfo
			recallP.Conn = p.Conn
			recallP.CommandID = 100
			recallP.Id = p.Id
			recallP.RegisterID = registerID

			recallJsonBuf, err := pkg.JsonMarshal(recallP)
			if err != nil {
				goLog.Error("recall json marshal error ", err)
				continue
			}

			func() {
				pipeLock.Lock()
				defer pipeLock.Unlock()
				pipeprotocol.SendMessage(pipeClient, recallJsonBuf)
				goLog.Debug("send recall info ", p.Id)
			}()
			continue
		}

		goLog.Debug("decode onec keys: ", p.Counts)
		aesOnec := pipeprotocol.AesNewNonece(p.Counts)
		aesDeBuf, err := pipeprotocol.AesDecode(aead, aesOnec, p.Buf)
		if err != nil {
			goLog.Error("aes decode error", err)
			continue
		}

		sendLen := 0
		for sendLen < len(aesDeBuf) {
			n, err := cInfo.Conn.Write(aesDeBuf[sendLen:])
			if err != nil {
				goLog.Error("write real client error ", err)
				break
			}
			sendLen += n
		}
	}
}
