package kpeering

import (
	"context"
	"net"
	"sync"

	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
)

type Kpeering struct {
	Frontdoor model.Frontdoor
	Backdoor  model.Backdoor
}

func (cfg *Kpeering) Start() {
	ctx := context.Background()
	reqChan := make(chan []byte)
	resChan := make(chan []byte)

	go cfg.startBackdoorTCPServer(reqChan, resChan)
	go cfg.startFrontdoorTCPServer(reqChan, resChan)

	<-ctx.Done()
}

/*
req -> frontdoorConn -> reqChan -> backdoorConn
- 										|
res <- frontdoorConn <- resChan <- backdoorConn
*/
func (cfg *Kpeering) startBackdoorTCPServer(reqChan chan []byte, resChan chan []byte) {
	var conn net.Conn
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			wg.Wait()

			err := io.ReadTo(conn, resChan)

			logger.Z.Errorf("[backdoor] connection is closed, %v", err)

			mutex.Lock()
			conn.Close()
			conn = nil
			wg.Add(1)

			// notify write goroutine to close
			msg := []byte("close")
			reqChan <- msg

			mutex.Unlock()

			logger.Z.Infoln("[backdoor] wait for next connection")
		}
	}()

	go func() {
		for {
			wg.Wait()

			err := io.WriteTo(reqChan, conn)

			logger.Z.Errorf("[backdoor] write coroutine is close as well, %v", err)
		}
	}()

	io.StartTCPServer(cfg.Backdoor.Address(),
		func(s string) {
			logger.Z.Infof("[backdoor] server is started on %s", s)
		},
		func(c net.Conn) {
			mutex.Lock()
			if conn == nil {
				logger.Z.Infoln("[backdoor] New connection from", conn.RemoteAddr())
				conn = c
				wg.Done()
			} else {
				logger.Z.Errorf("[backdoor] connection is already exists")
				c.Close()
			}
			mutex.Unlock()
		},
	)
}

func (cfg *Kpeering) startFrontdoorTCPServer(reqChan chan []byte, resChan chan []byte) {
	var conn net.Conn
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			wg.Wait()

			err := io.ReadTo(conn, reqChan)

			logger.Z.Errorf("[frontdoor] connection is closed, %v", err)

			mutex.Lock()
			conn.Close()
			conn = nil
			wg.Add(1)

			// notify write goroutine to close
			msg := []byte("close")
			resChan <- msg

			mutex.Unlock()

			logger.Z.Infoln("[frontdoor] wait for next connection")
		}
	}()

	go func() {
		for {
			wg.Wait()

			err := io.WriteTo(resChan, conn)

			logger.Z.Errorf("[frontdoor] write coroutine is close as well, %v", err)
		}
	}()

	io.StartTCPServer(cfg.Frontdoor.Address(),
		func(s string) {
			logger.Z.Infof("[frontdoor] server is started on %s", s)
		},
		func(c net.Conn) {
			mutex.Lock()
			if conn == nil {
				logger.Z.Infoln("[frontdoor] New connection from", conn.RemoteAddr())
				conn = c
				wg.Done()
			} else {
				logger.Z.Errorf("[frontdoor] connection is already exists")
				c.Close()
			}
			mutex.Unlock()
		},
	)
}
