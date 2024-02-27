package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ethInterfacingTestSuite struct {
	suite.Suite
}

func (suite *ethInterfacingTestSuite) SetupTest() {
	// setup test
}

func (suite *ethInterfacingTestSuite) TestIPControls() {
	err, stroutput := Setup_EthInterface()
	require.NoError(suite.T(), err, "Refresh_nmcli() should return nil on succesfull")
	assert.NotEqual(suite.T(), "", stroutput, "Refresh_nmcli() should return a non-empty string")
}

func TestEthInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(ethInterfacingTestSuite))
}
