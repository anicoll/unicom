package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/anicoll/unicom/internal/server"
)

type ServerUnitTestSuite struct {
	suite.Suite
	svc *server.Server
	tc  *mockTemporalClient
	db  *mockDatabase
}

func TestServerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ServerUnitTestSuite))
}

func (s *ServerUnitTestSuite) SetupSuite() {}

func (suite *ServerUnitTestSuite) SetupTest() {
	suite.tc = &mockTemporalClient{}
	suite.db = &mockDatabase{}
	suite.svc = server.New(zap.NewNop(), suite.tc, suite.db)
}

func (suite *ServerUnitTestSuite) AfterTest() {
	suite.tc.AssertExpectations(suite.T())
	suite.db.AssertExpectations(suite.T())
}
