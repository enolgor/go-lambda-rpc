package internal

import (
	"encoding/gob"
	"encoding/json"
	"io"
)

func EncodeJson(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}
func DecodeJson(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	return dec.Decode(r)
}

func EncodeGob(w io.Writer, v any) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(v)
}
func DecodeGob(r io.Reader, v any) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(r)
}
