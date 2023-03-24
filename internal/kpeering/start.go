package kpeering

import (
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

func (peering *Kpeering) Start() {
	var frontdoorConn net.Conn
	var backdoorConn net.Conn
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go io.StartTCPServer(peering.Frontdoor.Address(),
		func(s string) {
			logger.Z.Infof("frontdoor server is started on %s", s)
		},
		func(c net.Conn) {
			mutex.Lock()
			if frontdoorConn == nil {
				frontdoorConn = c
				wg.Done()
			} else {
				logger.Z.Errorf("frontdoor connection is already exists")
				c.Close()
			}
			mutex.Unlock()
		},
	)

	go io.StartTCPServer(peering.Backdoor.Address(),
		func(s string) {
			logger.Z.Infof("backdoor server is started on %s", s)
		},
		func(c net.Conn) {
			mutex.Lock()
			if backdoorConn == nil {
				backdoorConn = c
				wg.Done()
			} else {
				logger.Z.Errorf("backdoor connection is already exists")
				c.Close()
			}
			mutex.Unlock()
		},
	)

	for {
		wg.Wait()
		logger.Z.Infof("frontdoor and backdoor connections are ready")

		// loop until one of connections is closed
		err := io.BiFoward(frontdoorConn, backdoorConn)
		switch err {
		case io.ErrSourceDisconnected:
			logger.Z.Infof("frontdoor connection is closed")
			frontdoorConn.Close()
			frontdoorConn = nil
			wg.Add(1)
		case io.ErrTargetDisconnected:
			logger.Z.Infof("backdoor connection is closed")
			backdoorConn.Close()
			backdoorConn = nil
			wg.Add(1)
		}

		logger.Z.Infoln("wait for reconnection...")
	}
}
