package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/enolgor/go-lambda-rpc/internal"
	"github.com/enolgor/go-lambda-rpc/model"
)

type RPCHandler interface {
	Register(string, func(context.Context, Decode) (any, error))
	Handler(context.Context, *model.RPCEventInput) (*model.RPCEventOutput, error)
}

type Decode func(any) error

type rpcHandler struct {
	funcs   map[string]func(context.Context, Decode) (any, error)
	encoder func(io.Writer, any) error
	decoder func(io.Reader, any) error
}

func NewJsonHandler() RPCHandler {
	return &rpcHandler{
		funcs:   make(map[string]func(context.Context, Decode) (any, error)),
		encoder: internal.EncodeJson,
		decoder: internal.DecodeJson,
	}
}

func NewGobHandler() RPCHandler {
	return &rpcHandler{
		funcs:   make(map[string]func(context.Context, Decode) (any, error)),
		encoder: internal.EncodeGob,
		decoder: internal.DecodeGob,
	}
}

func (rh *rpcHandler) Register(path string, f func(context.Context, Decode) (any, error)) {
	rh.funcs[path] = f
}

func (rh *rpcHandler) Handler(ctx context.Context, input *model.RPCEventInput) (output *model.RPCEventOutput, err error) {
	output = &model.RPCEventOutput{}
	f, ok := rh.funcs[input.Path]
	if !ok {
		err = fmt.Errorf("path %s is not registered", input.Path)
		return
	}
	buffer := bytes.NewBuffer(input.Args)
	var result any
	if result, err = f(ctx, func(value any) error {
		return rh.decoder(buffer, value)
	}); err != nil {
		return
	}
	buffer = &bytes.Buffer{}
	if err = rh.encoder(buffer, result); err != nil {
		return
	}
	output.Response = buffer.Bytes()
	return
}
