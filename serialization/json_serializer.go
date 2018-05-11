package serialization

import (
	"encoding/json"
	"io"
)

type JSONSerializer struct{}

func JSON() *JSONSerializer {
	return &JSONSerializer{}
}

func (this *JSONSerializer) Serialize(source interface{}, destination io.Writer) error {
	return json.NewEncoder(destination).Encode(source)
}
func (this *JSONSerializer) Deserialize(source io.Reader, destination interface{}) error {
	return json.NewDecoder(source).Decode(destination)
}
