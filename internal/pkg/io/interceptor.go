package io

import (
	"context"
	"net"
	"sync"

	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
)

/*
client (conn) -> TCPInterceptor -> requestChan
client (conn) <- TCPInterceptor <- responseChan
*/
type TCPInterceptor struct {
	ctx       context.Context
	mutex     sync.Mutex
	wg        sync.WaitGroup
	conn      net.Conn
	address   string
	readInto  chan []byte
	writeFrom chan []byte
}

func NewTCPInterceptor(
	ctx context.Context,
	cfg model.Interceptor,
	requestChan chan []byte,
	responseChan chan []byte,
) *TCPInterceptor {
	return &TCPInterceptor{
		ctx:       context.WithValue(ctx, keyComponentID, cfg.Name),
		mutex:     sync.Mutex{},
		wg:        sync.WaitGroup{},
		address:   cfg.Address(),
		readInto:  requestChan,
		writeFrom: responseChan,
	}
}

func (t *TCPInterceptor) Run() {
	t.wg.Add(1)

	// read data and put into the forward channel
	go func() {
		for {
			t.wg.Wait()

			err := ReadTo(t.conn, t.readInto)

			t.close()

			logger.Z.Errorf("[%s] connection is closed, %v", t.name(), err)
			logger.Z.Infof("[%s] wait for next connection", t.name())
		}
	}()

	// write backword channel data to the connection
	go func() {
		for {
			t.wg.Wait()

			err := WriteTo(t.writeFrom, t.conn)

			logger.Z.Errorf("[%s] write coroutine is close as well, %v", t.name(), err)
		}
	}()

	StartTCPServer(
		t.address,
		func(s string) {
			logger.Z.Infof("[%s] server is started on %s", t.name(), s)
		},
		t.newConnection,
	)
}

func (t *TCPInterceptor) name() string {
	return t.ctx.Value(keyComponentID).(string)
}

func (t *TCPInterceptor) close() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.wg.Add(1)
	t.conn.Close()
	t.conn = nil

	// notify write goroutine to close
	msg := []byte("close")
	t.writeFrom <- msg
}

func (t *TCPInterceptor) newConnection(c net.Conn) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.conn == nil {
		logger.Z.Infof("[%s] New connection from %s", t.name(), c.RemoteAddr())
		t.conn = c
		t.wg.Done()
	} else {
		logger.Z.Errorf("[%s] connection is already exists", t.name())
		c.Close()
	}
}
