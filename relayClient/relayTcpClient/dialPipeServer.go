package relaytcpclient

import (
	"crypto/cipher"
	"fmt"
	dbcache "goRelay/dbCache"
	pipeprotocol "goRelay/pipeProtocol"
	"goRelay/pkg"
	"net"
	"strconv"
	"time"
)

func ConnectToPipeServer(addr string, registerID string, cache *dbcache.Caches) {

	clientAEAD = make(map[string]cipher.AEAD, 0)

	for {
		pipeClientConn, err = net.Dial("tcp", addr)
		if err != nil {
			goLog.Error("Dial relay server error", err)
			time.Sleep(time.Second * time.Duration(sleepTimeSec))
			continue
		}
		goLog.Debug("to ", pipeClientConn.RemoteAddr().String(), " send relay Conn information")

		relayConn = pipeClientConn

		var p pipeprotocol.ClientProtocolInfo
		p.RegisterID = registerID
		p.Buf = []byte("relayConn")

		relayConnBuf, err := pkg.JsonMarshal(p)
		if err != nil {
			goLog.Error("relay conn is error")
		}
		pipeprotocol.SendMessage(pipeClientConn, relayConnBuf)

		pipeResponse, err := pipeprotocol.RecvMessgae(pipeClientConn)
		if err != nil {
			goLog.Error("No available connections on the relay server. error: ", err)
			time.Sleep(time.Second * time.Duration(sleepTimeSec))
			continue
		}

		var replyRelayConnBuf pipeprotocol.ClientProtocolInfo
		pkg.JsonUnmarshal(pipeResponse, &replyRelayConnBuf)

		if string(replyRelayConnBuf.Buf) == "isok" {
			goLog.Info(pipeClientConn.LocalAddr().String(), " is Relay server address , connection established successfully.")
		}

		for {
			msg, err := pipeprotocol.RecvMessgae(pipeClientConn)
			if err != nil {
				goLog.Error(pipeClientConn.RemoteAddr().Network(), " Relay server connection is error. err: ", err)
				break
			}
			var p pipeprotocol.ClientProtocolInfo
			err = pkg.JsonUnmarshal(msg, &p)
			if err != nil {
				goLog.Error("json unmarshal error", err)
			}

			_, aeadIsok := clientAEAD[p.Id]
			if !aeadIsok {
				goLog.Debug("new aes cipher id: ", p.Id)
				aead, err := pipeprotocol.AesNewCipher(p.Id)
				if err != nil {
					goLog.Error("aes new cipher error ", err)
					continue
				}

				clientAEAD[p.Id] = aead
			}

			goLog.Debug("from pipe server read messages, id: ", p.Id, " data len: ", len(msg))

			isok := func(pConn string) bool {
				clientMapMutex.Lock()
				defer clientMapMutex.Unlock()

				_, isok := clientConnections[pConn]
				return isok
			}(p.Conn)
			if !isok {
				// realServerAddr, isok := realServerInfoMap[p.Id]
				stopClientID := fmt.Sprintf("Stop_%s", p.Id)
				if cache.Exists(stopClientID) {
					goLog.Debug("realservers is disabled ,id: ", p.Id)
					continue
				}

				realServerAddr, isok := cache.Get(p.Id)
				if !isok {
					goLog.Error("not fount realservers ,id: ", p.Id)
					continue
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

				AddTX(p.Id, len(p.Buf), cache)

				realServer, err := net.Dial("tcp", realServerAddr)
				if err != nil {
					goLog.Error("dial real server error", err)
					continue
				}

				var c pipeprotocol.ClientInfo
				c.Conn = realServer
				c.Time = time.Now().Unix()

				func() {
					clientMapMutex.Lock()
					defer clientMapMutex.Unlock()
					clientConnections[p.Conn] = c
				}()

				go connectionRealServer(realServer, p.Conn, p.Id, cache)

			}
			func() {
				cInfo, _ := clientConnections[p.Conn]

				if p.CommandID == 100 {
					goLog.Info("close ", p.Id, " connecttion, commandid: ", p.CommandID)
					cInfo.Conn.Close()
					return
				}

				AddTX(p.Id, len(p.Buf), cache)

				goLog.Debug("aes decode ", p.Counts)
				aesOnec := pipeprotocol.AesNewNonece(p.Counts)
				aesDeBuf, err := pipeprotocol.AesDecode(clientAEAD[p.Id], aesOnec, p.Buf)
				if err != nil {
					goLog.Error("aes decode error ", err)
					return
				}

				sendLen := 0
				for sendLen < len(aesDeBuf) {
					n, err := cInfo.Conn.Write(aesDeBuf[sendLen:])
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
			}()

			// func() {
			// 	sendMutex.Lock()
			// 	defer sendMutex.Unlock()

			// 	// if nil == clientConnection {
			// 	// 	goLog.Error("relayClient Conn is null")
			// 	// 	return
			// 	// }

			// 	goLog.Debug("from pipe server read messages, id: ", p.Id, " data len: ", len(msg))
			// 	// pipeprotocol.SendMessage(clientConnection, msg)

			// }()

		}
	}
}

func AddRX(id string, dateLen int, cache *dbcache.Caches) {
	var newRX string = fmt.Sprintf("Usage_RX_%s", id)
	strTmpRx, rxIsok := cache.Get(newRX)
	if rxIsok {
		tmp_rx, err := strconv.ParseUint(strTmpRx, 10, 64)
		if err != nil {
			goLog.Error("rx to uint64 error", err)
			tmp_rx = 0

		}
		tmp_rx += uint64(dateLen)
		cache.Set(newRX, strconv.FormatUint(tmp_rx, 10))
	} else {
		cache.Set(newRX, strconv.Itoa(dateLen))
	}
}

func AddTX(id string, dateLen int, cache *dbcache.Caches) {
	var newTx string = fmt.Sprintf("Usage_TX_%s", id)

	strTmpTx, rxIsok := cache.Get(newTx)
	if rxIsok {
		tmp_ex, err := strconv.ParseUint(strTmpTx, 10, 64)
		if err != nil {
			goLog.Error("rx to uint64 error", err)
			tmp_ex = 0
		}
		tmp_ex += uint64(dateLen)
		cache.Set(newTx, strconv.FormatUint(tmp_ex, 10))
	} else {
		cache.Set(newTx, strconv.Itoa(dateLen))
	}
}
