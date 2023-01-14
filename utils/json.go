package utils

import (
	"bytes"
	"encoding/json"
)

func JsonUnEscape(v interface{}) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	return w.Bytes(), err
}

func CovertByJson(src, dest interface{}) error {
	bts, err := json.Marshal(src)
	if err == nil {
		err = json.Unmarshal(bts, dest)
	}
	return err
}

func MustJson(v interface{}) []byte {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bs
}

func MustJsonUnEscape(v interface{}) []byte {
	w := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
	return w.Bytes()
}

func MustJsonString(v interface{}) string {
	return string(MustJson(v))
}
