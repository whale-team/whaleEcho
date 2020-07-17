package handler_test

import (
	"context"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
	"google.golang.org/protobuf/proto"
)

func NewSuite(host, port string) *wsSuite {
	return &wsSuite{
		host:     host,
		port:     port,
		Ctx:      context.Background(),
		wsDialer: ws.DefaultDialer,
		T:        ginkgo.GinkgoT(),
	}
}

type wsSuite struct {
	wsDialer ws.Dialer
	host     string
	port     string
	server   *wsserver.SocketServer
	Ctx      context.Context
	T        ginkgo.GinkgoTInterface
}

func (s *wsSuite) setupServer(handler wshandler.Handler) {
	s.server = wsserver.New()
	s.server.Handler = logMiddleware(handler.Routing)
	s.server.ErrHandler = func(err error) {
		log.Error().Stack().Err(err).Msg("read message failed")
	}
	s.server.ConnBuildHandleFunc = wsserver.ConnBuildHandle
	s.server.ConnCloseHandler = func(c *wsserver.Context) error {
		log.Debug().Msg("conn closed")
		return nil
	}
}

func (s wsSuite) runServer() {
	go s.server.ListenAndServe(s.host, s.port)
}

func (s wsSuite) Dial() (net.Conn, error) {
	conn, _, _, err := s.wsDialer.Dial(s.Ctx, "ws://"+s.host+":"+s.port)
	return conn, err
}

func (w wsSuite) SendCommand(conn net.Conn, command *echoproto.Command) error {
	data, err := proto.Marshal(command)
	if err != nil {
		return err
	}
	return wsutil.WriteClientBinary(conn, data)
}

func (w wsSuite) ReadResp(conn net.Conn) (*echoproto.CommandResp, error) {
	data, _, err := wsutil.ReadServerData(conn)
	if err != nil {
		return nil, err
	}
	resp := &echoproto.CommandResp{}

	if err := proto.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func logMiddleware(handler wsserver.HandleFunc) wsserver.HandleFunc {
	return func(c *wsserver.Context) error {
		log.Debug().Msg("rec message")
		return handler(c)
	}
}
