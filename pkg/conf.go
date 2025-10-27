package pkg

import (
	"errors"
	"fmt"
	"os"
)

var (
	LogLevel = InfoLevel
)

func LoadConfig(filename string, v interface{}) error {
	fbody, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("read configure file ", filename, "error", err)
		return errors.New("read config error")
	}

	err = YamlUnmarshal(fbody, v)
	if err != nil {
		fmt.Println("json unmarshal error", err)
		return errors.New("json unmarshal error")
	}

	return nil
}
