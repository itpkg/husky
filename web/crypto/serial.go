package crypto

import (
	"bytes"
	"encoding/gob"
)

//Serial serial
type Serial struct {
}

//From object from bytes
func (p *Serial) From(b []byte, v interface{}) error {
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(b)
	return dec.Decode(v)
}

//To object to bytes
func (p *Serial) To(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err == nil {
		return buf.Bytes(), nil
	} else {
		return nil, err
	}
}
