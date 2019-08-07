package common

import (
	"bytes"
	"encoding/gob"
)
//从interface{}转换为bytes
func GetBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
//从bytes转换为interface{}
func GetInterface(bts []byte,data  interface{}) (error) {

	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
