package serialization

import "io"

type Serializer interface {
	Serialize(interface{}, io.Writer) error
	Deserialize(io.Reader, interface{}) error
}
