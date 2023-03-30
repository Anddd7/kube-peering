package connectors

import (
	"context"
	"sync"

	"github.com/kube-peering/internal/pkg/model"
)

/*
client (conn) <- TunnelServer <- requestChan
client (conn) -> TunnelServer -> responseChan
*/
func NewTunnelServer(
	ctx context.Context,
	cfg model.Tunnel,
	requestChan chan []byte,
	responseChan chan []byte,
) *TCPInterceptor {
	return &TCPInterceptor{
		ctx:       context.WithValue(ctx, keyComponentID, cfg.Name),
		mutex:     sync.Mutex{},
		wg:        sync.WaitGroup{},
		address:   cfg.Address(),
		readInto:  responseChan,
		writeFrom: requestChan,
	}
}

/*
responseChan -> TunnelClient -> server (conn)
requestChan <- TunnelClient <- server (conn)
*/
func NewTunnelClient(
	ctx context.Context,
	cfg model.Tunnel,
	requestChan chan []byte,
	responseChan chan []byte,
) *TCPForwarder {
	return &TCPForwarder{
		ctx:          context.WithValue(ctx, keyComponentID, cfg.Name),
		address:      cfg.Address(),
		forwardChan:  responseChan,
		backwordChan: requestChan,
	}
}
