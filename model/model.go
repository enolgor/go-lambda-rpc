package model

type RPCEventInput struct {
	Path string `json:"path"`
	Args []byte `json:"args"`
}

type RPCEventOutput struct {
	Response []byte `json:"response"`
}

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}
