package pipeprotocol

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

func dataCompressForGzip(content []byte) []byte {
	var data bytes.Buffer
	gzipWrite, err := gzip.NewWriterLevel(&data, gzip.BestCompression)
	if err != nil {
		return nil
	}

	writeLen := 0
	for writeLen < len(content) {
		n, err := gzipWrite.Write(content[writeLen:])
		if err != nil {
			return nil
		}
		writeLen += n
	}

	err = gzipWrite.Flush()
	if err != nil {
		return nil
	}
	err = gzipWrite.Close()
	if err != nil {
		return nil
	}
	return data.Bytes()
}

func dataDecompressForGzip(compressedData []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		return nil, errors.New(fmt.Sprintln("data decompress error: ", err))
	}

	content, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("io read all error: ", err))
	}

	return content, nil
}
