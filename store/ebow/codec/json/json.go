package json

import (
	"encoding/json"

	"at.ourproject/energystore/store/ebow/codec"
)

type Codec struct{}

func (c Codec) Marshal(v interface{}, in []byte) (out []byte, err error) {
	return json.Marshal(v)
}

func (c Codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (c Codec) Format() codec.Format {
	return codec.JSON
}
