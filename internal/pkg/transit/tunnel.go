package transit

type Tunnel interface {
	Start()
	Send()
}
