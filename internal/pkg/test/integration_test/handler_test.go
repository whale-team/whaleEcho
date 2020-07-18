package handler_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vicxu416/goinfra/zlogging"
	"github.com/whale-team/whaleEcho/configs"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker/natsbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"github.com/whale-team/whaleEcho/pkg/natspool"
	"go.uber.org/fx"
)

func TestHandler(t *testing.T) {
	if err := setupSuite(); err != nil {
		t.Fatalf("setup suite failed, err:%+v", err)
		t.FailNow()
	}
	time.Sleep(50 * time.Millisecond)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Websocket Handler Spec")
}

var suite *wsSuite

func relativePath(path string) string {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	return filepath.Join(dir, path)
}

func setupSuite() error {
	suite = NewSuite("127.0.0.1", "22222")

	if err := suite.initData(); err != nil {
		return err
	}

	os.Setenv("CONFIG_NAME", "app-test")
	config, err := configs.InitConfiguration()
	if err != nil {
		return err
	}

	zlogging.SetupLogger(config.Log)

	handler := wshandler.Handler{}

	app := fx.New(
		fx.Supply(config.Nats),
		fx.Provide(natspool.NewClient, roomscenter.New),
		fx.Provide(natsbroker.New, service.New, wshandler.New),
		fx.Populate(&handler, &suite.broker, &suite.center),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		return err
	}
	if err := app.Stop(ctx); err != nil {
		return err
	}
	suite.setupServer(handler)
	suite.runServer()

	return nil
}
