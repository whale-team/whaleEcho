package handler_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
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

func setupSuite() error {
	suite = NewSuite("127.0.0.1", "22222")

	handler := wshandler.Handler{}
	handler.SetupRoutes()

	suite.setupServer(handler)
	suite.runServer()

	return nil
}
