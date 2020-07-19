package testutil

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

// NewSuite construct testsuite
func NewSuite() *TestSuite {
	return &TestSuite{
		Ctx:        context.Background(),
		GrpcTester: &GrpcTester{},
		DataTester: &DataTester{TestData: make(map[string][]byte)},
		T:          &asserter{},
	}
}

// TestSuite composite all tester
type TestSuite struct {
	Ctx context.Context
	app *fx.App
	*GrpcTester
	*DataTester
	T *asserter
}

func (suite TestSuite) Err() error {
	return suite.T.err
}

func (suite *TestSuite) ClearErr() {
	suite.T.err = nil
}

// RunApp start fx App
func (suite *TestSuite) RunApp(app *fx.App) error {
	suite.app = app

	return suite.app.Start(suite.Ctx)
}

// Done release all resources
func (suite *TestSuite) Done() error {
	if err := suite.app.Stop(suite.Ctx); err != nil {
		return err
	}

	if err := suite.GRPConn.Close(); err != nil {
		return err
	}

	if err := suite.DB.Close(); err != nil {
		return err
	}

	return nil
}

type asserter struct {
	err error
}

func (a *asserter) Errorf(format string, args ...interface{}) {
	if a.err != nil {
		a.err = fmt.Errorf(format+a.err.Error(), args...)
	} else {
		a.err = fmt.Errorf(format, args...)
	}
}
