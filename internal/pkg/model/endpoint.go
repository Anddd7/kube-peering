package model

import (
	"fmt"
)

type Endpoint struct {
	Name       string
	Protocol   string
	Host       string
	ListenPort int
}

func (e Endpoint) Address() string {
	return fmt.Sprintf("%s:%d", e.Host, e.ListenPort)
}

func (e Endpoint) IsTCP() bool {
	return e.Protocol == "tcp"
}

func (e Endpoint) IsHTTP() bool {
	return e.Protocol == "http"
}

func (e Endpoint) IsHTTPS() bool {
	return e.Protocol == "https"
}
