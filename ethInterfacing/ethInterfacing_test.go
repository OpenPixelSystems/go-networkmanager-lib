package ethInterfacing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ethInterfacingTestSuite struct {
	suite.Suite
}

// Test SetIPAddr
// - err should be nil on succesfull
func (suite *ethInterfacingTestSuite) TestSetIPAddr() {
	err := SetIPAddr()
	assert.Equal(suite.T(), nil, err, "SetIPAddr() should return nil on succesfull")
}

// Test get_original_interface_setting
// - err should be nil on succesfull
func (suite *ethInterfacingTestSuite) TestGetInterfaceSettings() {
	err, ethface, ethaddr, ethgateway := Get_interface_settings()
	assert.Equal(suite.T(), nil, err, "get_original_interface_setting() should return nil on succesfull")
	assert.Equal(suite.T(), yourInterfaceName, ethface, "get_original_interface_setting() should return yourInterfaceName")
	assert.Equal(suite.T(), ip_addr, ethaddr, "get_original_interface_setting() should return ip_addr")
	assert.Equal(suite.T(), defgateway, ethgateway, "get_original_interface_setting() should return gateway")
}

// Test SetIPMode
// - err should be nil on succesfull
func (suite *ethInterfacingTestSuite) TestSetIpMode() {
	err := SetIPMode()
	assert.Equal(suite.T(), nil, err, "SetIpMode() should return nil on succesfull")
}

// Test Refresh_nmcli
func (suite *ethInterfacingTestSuite) TestRefreshNmcli() {
	err, stroutput := Refresh_nmcli()
	assert.Equal(suite.T(), nil, err, "Refresh_nmcli() should return nil on succesfull")
	assert.NotEqual(suite.T(), "", stroutput, "Refresh_nmcli() should return a non-empty string")
	assert.Contains(suite.T(), stroutput, ip_addr, "Refresh_nmcli() should return a string containing ip_addr")
	assert.Contains(suite.T(), stroutput, yourInterfaceName, "Refresh_nmcli() should return a string containing yourInterfaceName")
}

// Test SetDefaultGateway
// - err should be nil on succesfull
func (suite *ethInterfacingTestSuite) TestSetDefaultGateway() {
	err := SetDefaultGateway()
	assert.Equal(suite.T(), nil, err, "SetDefaultGateway() should return nil on succesfull")
}

func TestEthInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(ethInterfacingTestSuite))
}
