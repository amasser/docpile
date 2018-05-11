package serialization

import (
	"encoding/json"
	"io"
)

type Default struct{}

func New() *Default {
	return &Default{}
}

func (this *Default) Serialize(source interface{}, destination io.Writer) error {
	return json.NewEncoder(destination).Encode(source)
}
func (this *Default) Deserialize(source io.Reader, destination interface{}) error {
	return json.NewDecoder(source).Decode(destination)
}
