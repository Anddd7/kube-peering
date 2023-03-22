package kpeering

import (
	"net"

	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/model"
)

type Kpeering struct {
	Frontdoor model.Frontdoor
	Backdoor  model.Backdoor
}

func (peering *Kpeering) Start() {
	frontdoorConnChan := make(chan net.Conn)
	backdoorConnChan := make(chan net.Conn)

	go io.AcceptConnections(peering.Frontdoor.Name, peering.Frontdoor.Protocol, peering.Frontdoor.Address(), frontdoorConnChan)
	go io.AcceptConnections(peering.Backdoor.Name, "tcp", peering.Backdoor.Address(), backdoorConnChan)

	for {
		select {
		case f := <-frontdoorConnChan:
			b := <-backdoorConnChan
			go io.BiFoward(peering.Frontdoor.Name, f, peering.Backdoor.Name, b)
		case b := <-backdoorConnChan:
			f := <-frontdoorConnChan
			go io.BiFoward(peering.Frontdoor.Name, f, peering.Backdoor.Name, b)
		}
	}
}
