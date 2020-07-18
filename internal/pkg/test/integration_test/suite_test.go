package handler_test

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/onsi/ginkgo"

	"github.com/rs/zerolog/log"
	"github.com/vicxu416/goinfra/testutil"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
	"google.golang.org/protobuf/proto"
)

func NewSuite(host, port string) *wsSuite {
	return &wsSuite{
		TestSuite: testutil.NewSuite(),
		host:      host,
		port:      port,
		Ctx:       context.Background(),
		wsDialer:  websocket.DefaultDialer,
		T:         ginkgo.GinkgoT(),
		rooms:     make([]*echoproto.Room, 0, 1),
		users:     make([]*echoproto.User, 0, 1),
		messages:  make([]*echoproto.Message, 0, 1),
	}
}

type wsSuite struct {
	*testutil.TestSuite
	wsDialer *websocket.Dialer
	host     string
	port     string
	server   *wsserver.SocketServer
	broker   msgbroker.MsgBroker
	Ctx      context.Context
	T        ginkgo.GinkgoTInterface
	rooms    []*echoproto.Room
	users    []*echoproto.User
	messages []*echoproto.Message
	center   *roomscenter.Center
}

func (s *wsSuite) setupServer(handler wshandler.Handler) {
	s.server = wsserver.New()
	s.server.Handler = handler.Handle
	s.server.ErrHandler = func(c *wsserver.Context, err error) {
		log.Error().Stack().Err(err).Msg("read message failed")

	}
	s.server.ConnBuildHandleFunc = wsserver.ConnBuildHandle
	s.server.ConnCloseHandler = func(c *wsserver.Context) error {
		log.Debug().Msg("conn closed")
		c.Close()
		return nil
	}
}

func (s *wsSuite) initData() error {
	if err := suite.LoadTestData("../testdata"); err != nil {
		return err
	}
	if err := suite.UnmarshalTestData("rooms", &suite.rooms); err != nil {
		return err
	}
	if err := suite.UnmarshalTestData("users", &suite.users); err != nil {
		return err
	}
	if err := suite.UnmarshalTestData("messages", &suite.messages); err != nil {
		return err
	}
	return nil
}

func (s wsSuite) runServer() {
	go s.server.ListenAndServe(s.host, s.port)
}

func (s wsSuite) Dial() (*websocket.Conn, error) {
	conn, _, err := s.wsDialer.Dial("ws://"+s.host+":"+s.port, nil)
	return conn, err
}

func (w wsSuite) SendCommand(conn *websocket.Conn, command *echoproto.Command) error {
	data, err := proto.Marshal(command)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.BinaryMessage, data)
}

func (w wsSuite) ReadResp(conn *websocket.Conn) (*echoproto.Message, error) {
	_, data, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	resp := &echoproto.Message{}

	if err := proto.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
