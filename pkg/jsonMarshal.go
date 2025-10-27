package pkg

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

func JsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func JsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, &v)
}

func YamlUnmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
