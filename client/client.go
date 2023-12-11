package client

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/enolgor/go-lambda-rpc/internal"
	"github.com/enolgor/go-lambda-rpc/model"
)

type Client interface {
	Call(path string, args, result any) (err error)
}

type client struct {
	srv          *lambda.Client
	functionName string
	encoder      func(io.Writer, any) error
	decoder      func(io.Reader, any) error
}

func NewJsonClient(conf aws.Config, functionName string) Client {
	return &client{
		srv:          lambda.NewFromConfig(conf),
		functionName: functionName,
		encoder:      internal.EncodeJson,
		decoder:      internal.DecodeJson,
	}
}

func NewGobClient(conf aws.Config, functionName string) Client {
	return &client{
		srv:          lambda.NewFromConfig(conf),
		functionName: functionName,
		encoder:      internal.EncodeGob,
		decoder:      internal.DecodeGob,
	}
}

func (cli *client) Call(path string, args, result any) (err error) {
	input := &model.RPCEventInput{}
	input.Path = path
	buffer := &bytes.Buffer{}
	if err = cli.encoder(buffer, args); err != nil {
		return
	}
	input.Args = buffer.Bytes()
	output, err := cli.srv.Invoke(context.Background(), &lambda.InvokeInput{
		FunctionName: aws.String(cli.functionName),
		Payload:      buffer.Bytes(),
	})
	if err != nil {
		return err
	}
	if output.FunctionError != nil {
		err = errors.New(*output.FunctionError)
		return
	}
	buffer = bytes.NewBuffer(output.Payload)
	return cli.decoder(buffer, result)
}
