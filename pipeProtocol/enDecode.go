package pipeprotocol

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"goRelay/pkg"
)

func xor(buf []byte, s string) []byte {
	xs := make([]byte, 0)
	for i := 0; i < len(buf); i++ {
		xs = append(xs, buf[i]^(s[i%len(s)]))
	}
	return xs
}

func enDeCode(buf []byte, types int) []byte {
	if types == 1 {
		// 加密
		xbuf1 := xor(buf, Keys[0])
		return xbuf1
	} else if types == 2 {
		// 解密
		xbuf1 := xor(buf, Keys[0])
		return xbuf1
	}
	return nil
}

func Encode(s []byte) []byte {
	return enDeCode(s, 1)
}

func Decode(s []byte) []byte {
	return enDeCode(s, 2)
}

func AesNewCipher(id string) (cipher.AEAD, error) {
	const appstring = "zsNAmq0WW9LVgXCqeR9I7IPa9OomaVtc"
	keyLen := 32

	key := fmt.Sprintf("%s%s", id, appstring)

	newKey := make([]byte, 0)

	newKey = append(newKey, []byte(key[:keyLen])...)

	block, err := aes.NewCipher([]byte(newKey))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm, nil
}

func AesNewNonece(id string) []byte {
	const noncestring = "59wbagDTZGIx"

	return []byte(pkg.IDHash(fmt.Sprintf("%s%s", id, noncestring))[:12])
}

func AesEncode(aesgcm cipher.AEAD, nonce []byte, data []byte) []byte {
	return aesgcm.Seal(nil, nonce, data, nil)
}

func AesDecode(aesgcm cipher.AEAD, nonce []byte, data []byte) ([]byte, error) {
	return aesgcm.Open(nil, nonce, data, nil)
}
