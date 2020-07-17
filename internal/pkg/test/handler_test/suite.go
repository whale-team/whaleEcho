package handler

import "github.com/gobwas/ws"

type wsSuite struct {
	wsDialer ws.Dialer
	host     string
}

func (s wsSuite) runServer() {
}
