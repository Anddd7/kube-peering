package transit

type TransitRouter struct {
	// intercept the traffic and send to the tunnel or proxy
	interceptror Interceptor
	// p2p tunnel with another transit router
	tunnel Tunnel
	// proxy for non-tunnel traffic
	proxy Proxy
	// route table for in/out bound
	routeTable RouteTable
}

// TODO feature: route table for in/out bound
type RouteTable struct {
}
