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
