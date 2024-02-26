package ethInterfacing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ethInterfacingTestSuite struct {
	suite.suite
}

// Test SetIPAddr
// - err should be nil on succesfull
func (suite *ethInterfacingTest) TestSetIPAddr()
{
	err := SetIPAddr()
	assert.equal(suite.T(), err, nil,  "SetIPAddr() should return nil on succesfull")
}

// Test get_original_interface_setting
// - err should be nil on succesfull
func (suite *ethInterfacingTest) TestGetOriginalInterfaceSetting()
{
	err := get_original_interface_setting()
	assert.equal(suite.T(), err, nil,  "get_original_interface_setting() should return nil on succesfull")
}	

// Test SetIPMode
// - err should be nil on succesfull
func (suite *ethInterfacingTest) TestSetIpMode()
{
	err := SetIpMode()
	assert.equal(suite.T(), err, nil,  "SetIpMode() should return nil on succesfull")
}

func TestEthInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(ethInterfacingTestSuite))
}