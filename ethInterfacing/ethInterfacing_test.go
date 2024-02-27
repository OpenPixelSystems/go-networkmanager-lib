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
func (suite *ethInterfacingTestSuite) TestGetOriginalInterfaceSetting() {
	err := Get_original_interface_setting()
	assert.Equal(suite.T(), nil, err, "get_original_interface_setting() should return nil on succesfull")
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
	assert.Contains(suite.T(), stroutput, defgateway, "Refresh_nmcli() should return a string containing defgateway")
	assert.Contains(suite.T(), stroutput, yourInterfaceName, "Refresh_nmcli() should return a string containing yourInterfaceName")
}

func TestEthInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(ethInterfacingTestSuite))
}
