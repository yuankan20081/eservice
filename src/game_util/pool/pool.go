package pool

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"time"
)

type NewFunc func(cc *grpc.ClientConn) interface{}

type Pool struct {
	q                chan *grpc.ClientConn // cache queue
	target           string                // server address
	ctx              context.Context
	newServiceClient NewFunc
}

func New(ctx context.Context, size int, target string, fn NewFunc) *Pool {
	return &Pool{
		q:                make(chan *grpc.ClientConn, size),
		target:           target,
		ctx:              ctx,
		newServiceClient: fn,
	}
}

func (p *Pool) pop() (*grpc.ClientConn, error) {
	select {
	case cc := <-p.q:
		return cc, nil
	default:
		ctx, _ := context.WithTimeout(p.ctx, time.Second*5)
		return grpc.DialContext(ctx, p.target, grpc.WithInsecure())
	}
}

func (p *Pool) push(cc *grpc.ClientConn) {
	// if pool is destroying
	select {
	case <-p.ctx.Done():
		cc.Close()
		return
	default:
		break
	}

	// if pool is full
	select {
	case p.q <- cc:
	default:
		cc.Close()
	}
}

func (p *Pool) Call(ctx context.Context, method string, request interface{}, opts ...grpc.CallOption) (interface{}, error) {
	cc, err := p.pop()
	if err != nil {
		log.Println("rpc client pool error:", err)
		return nil, err
	}

	iRpcClient := p.newServiceClient(cc)

	rpcClient := reflect.ValueOf(iRpcClient)
	rpcClientMethod := rpcClient.MethodByName(method)
	if rpcClientMethod.Kind() != reflect.Func {
		defer p.push(cc)
		message := fmt.Sprintf("%s is not a method of type %v", method, rpcClient.Type())
		return nil, errors.New(message)
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(request),
		// TODO: add opts...
	}

	resp := rpcClientMethod.Call(args)

	if len(resp) != 2 {
		defer p.push(cc)
		message := fmt.Sprintf("%s has too much return item", method)
		return nil, errors.New(message)
	}

	if err, ok := resp[1].Interface().(error); !ok && !resp[1].IsNil() {
		defer p.push(cc)
		message := fmt.Sprintf("second item return by %s is not nil or error", method)
		return nil, errors.New(message)
	} else {
		if err == grpc.ErrClientConnClosing {
			cc.Close()
			log.Println("*************")
			return nil, err
		}
		defer p.push(cc)
		return resp[0].Interface(), err
	}
}
