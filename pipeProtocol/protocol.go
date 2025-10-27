package pipeprotocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type ClientProtocolInfo struct {
	Id         string `json:"id"`
	CommandID  int    `json:"commandID"`
	Conn       string `json:conn`
	DirectConn string `json:"directConn"`
	Buf        []byte `json:"buf"`
	RegisterID string `json:"registerID"`
	Counts     string `json:"counts"`
}

type ClientInfo struct {
	Conn net.Conn
	Time int64
}

func SendMessage(conn net.Conn, message []byte) {

	enMsgs := Encode(message)

	msg := dataCompressForGzip(enMsgs)
	if msg == nil {
		return
	}

	msgLen := uint32(len(msg))

	headBuf := make([]byte, headerLen)
	binary.LittleEndian.PutUint32(headBuf, msgLen)

	sendLen := 0
	for sendLen < headerLen {
		n, err := conn.Write(headBuf[sendLen:])
		if err != nil {
			return
		}
		sendLen += n
	}

	sendLen = 0
	for sendLen < int(msgLen) {
		n, err := conn.Write([]byte(msg)[sendLen:])
		if err != nil {
			return
		}
		sendLen += n
	}
}

func RecvMessgae(conn net.Conn) ([]byte, error) {

	var msgLen uint32
	headBuf := make([]byte, headerLen)

	recvLen := 0
	for recvLen < headerLen {
		n, err := conn.Read(headBuf[recvLen:])
		if err != nil {
			if io.EOF == err {
				return nil, io.EOF
			}
			return nil, errors.New(fmt.Sprintf("read", conn.RemoteAddr().String(), " head buf error: ", err))
		}
		recvLen += n
	}

	msgLen = binary.LittleEndian.Uint32(headBuf)

	buf := make([]byte, msgLen)

	recvLen = 0
	for recvLen < int(msgLen) {
		n, err := conn.Read(buf[recvLen:])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("read", conn.RemoteAddr().String(), " data buf error: ", err))
		}
		recvLen += n
	}

	deCompressBuf, err := dataDecompressForGzip(buf[:])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("de compress buf error: ", err))
	}

	deBuf := Decode(deCompressBuf)

	return deBuf, nil
}
