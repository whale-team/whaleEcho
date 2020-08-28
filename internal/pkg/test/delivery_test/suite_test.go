package delivery_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/vicxu416/goinfra/testutil"
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/configs"
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/listener"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/dispatcher"
	"github.com/whale-team/whaleEcho/internal/pkg/repository/db"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/stanclient"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

var suite *testSuite

func TestDelivery(t *testing.T) {
	var err error
	suite, err = setupSuite("127.0.0.1", "13333")
	if err != nil {
		t.Fatalf("setup suite failed, err:%+v", err)
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Websocket Handler Spec")
	suite.teardown()
}

type testSuite struct {
	serv *wsserver.Server
	app  *fx.App
	ctx  context.Context
	T    GinkgoTInterface
	addr string
	port string
	stan *stanclient.Client
	*testutil.TestSuite
	repo        app.Repositorier
	rms         *dispatcher.Rooms
	redisClient *redis.Client

	rooms    []*echoproto.Room
	users    []*echoproto.User
	messages []*echoproto.Message
	buf      *bytes.Buffer
}

func dial(addr, port string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+addr+":"+port, nil)
	return conn, err
}

func (suite *testSuite) clear() {
	suite.buf.Reset()
	suite.rms.Clear()
	suite.redisClient.FlushAll(suite.ctx).Result()
}

func (suite *testSuite) setupServer() {
	go suite.serv.Start()
	time.Sleep(1 * time.Millisecond)
}

func (suite *testSuite) teardownServer() {
	suite.serv.Shutdown(1 * time.Millisecond)
}

func (suite *testSuite) initData() error {
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

func (suite testSuite) SendCommand(conn *websocket.Conn, command *echoproto.Command) error {
	data, err := proto.Marshal(command)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.BinaryMessage, data)
}

func (suite testSuite) ReadResp(conn *websocket.Conn) (*echoproto.Message, error) {
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

func setupSuite(addr, port string) (*testSuite, error) {
	os.Setenv("CONFIG_NAME", "app-test")
	config, err := configs.InitConfiguration()
	if err != nil {
		return nil, err
	}

	suite := testSuite{
		ctx:       context.Background(),
		T:         GinkgoT(),
		TestSuite: testutil.NewSuite(),
		buf:       new(bytes.Buffer),
		addr:      addr,
		port:      port,
	}

	writer := io.MultiWriter(suite.buf, os.Stdout)

	log.Logger = log.Output(writer).With().Logger()

	suite.serv = wsserver.NewDefault()
	suite.serv.Addr = addr
	suite.serv.Port = port

	suite.app = fx.New(
		fx.Supply(config, suite.serv),
		fx.Provide(stanclient.New, db.NewRedis, db.New, dispatcher.NewRooms, dispatcher.New),
		fx.Provide(service.New, listener.New, wshandler.New),
		fx.Invoke(wshandler.SetupHandler, listener.Listen),
		fx.Populate(&suite.stan, &suite.redisClient, &suite.rms, &suite.repo),
	)
	if err := suite.app.Start(suite.ctx); err != nil {
		return nil, err
	}
	if err := suite.initData(); err != nil {
		return nil, err
	}
	suite.setupServer()

	return &suite, nil
}

func (suite *testSuite) teardown() {
	suite.stan.Close()
}

func newCommand(msg proto.Message, typ echoproto.CommandType) (echoproto.Command, error) {
	com := echoproto.Command{}
	com.Type = typ
	data, err := proto.Marshal(msg)
	if err != nil {
		return com, err
	}
	com.Payload = data
	return com, nil
}
