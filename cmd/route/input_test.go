package route
  
import (
        "os"
        "testing"
	"github.com/stretchr/testify/assert"
)

func TestValidateAddress(t *testing.T) {
	assert := assert.New(t)
        assert.Equal(false, ValidateAddress("333.22.1.1"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(false, ValidateAddress("abc.22.1.1"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(false, ValidateAddress("22.abc.1.1"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(false, ValidateAddress("22.334.1.1"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(false, ValidateAddress("22.1.334.1"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(false, ValidateAddress("22.1.abc.1"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(false, ValidateAddress("22.1.1.abc"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(false, ValidateAddress("22.1.1.1111"), "Invalid Address, please follow IPv4 standard")
        assert.Equal(true,  ValidateAddress("22.1.1.1"), "Valid IP address")
}

func TestConvertPrefixLenToMask(t *testing.T) {
	assert := assert.New(t)
        assert.Equal("255.255.255.0", ConvertPrefixLenToMask("24"), "Converting Prefix lenth to right Mask in dotted format")
		
}

func TestValidatePrefixLen(t *testing.T) {
	assert := assert.New(t)
        assert.Equal(false, ValidatePrefixLen("40"), "Invalid Prefix length for IPv4. Prefix length should not be more than 32")
        assert.Equal(false, ValidatePrefixLen("abc"), "Invalid Prefix length for IPv4. Prefix length should not be string")
        assert.Equal(true, ValidatePrefixLen("22"), "Valid Prefix len")
}

func TestExtractNetworkAndPrefix(t *testing.T) {
	assert := assert.New(t)
        net, prefix := ExtractNetworkAndPrefix("10.10.10.10/24")
	assert.Equal(net, "10.10.10.10", "Extracted Network correctly")	
	assert.Equal(prefix, "24", "Extracted Network correctly")	
		
}

func TestInitPrefixSubnetTable(t *testing.T) {
	assert := assert.New(t)
        InitPrefixSubnetTable()
	assert.Equal(PrefixSubnetTable["25"],  "128", "Prefix Subnet Table initialized correctly")	
		
}

func TestInitializeNodeIP(t *testing.T) {
	assert := assert.New(t)
	input := new(Input)
	input.Network = "20.20.20.20"
	input.PrefixLen = "24"
	InitializeNodeIP(input)
	assert.Equal(input.NextAddress,  "20.20.20.0", "Calculated Next Address Correctly")	
	input.Network = "20.20.20.200"
	input.PrefixLen = "25"
	InitializeNodeIP(input)
	assert.Equal(input.NextAddress,  "20.20.20.128", "Calculated Next Address Correctly")	
	input.Network = "20.20.20.200"
	input.PrefixLen = "16"
	InitializeNodeIP(input)
	assert.Equal(input.NextAddress,  "20.20.0.0", "Calculated Next Address Correctly")	
	input.Network = "20.20.200.200"
	input.PrefixLen = "17"
	InitializeNodeIP(input)
	assert.Equal(input.NextAddress,  "20.20.128.0", "Calculated Next Address Correctly")	
	input.PrefixLen = "8"
	InitializeNodeIP(input)
	assert.Equal(input.NextAddress,  "20.0.0.0", "Calculated Next Address Correctly")	
	input.Network = "20.200.200.200"
	input.PrefixLen = "9"
	InitializeNodeIP(input)
	assert.Equal(input.NextAddress,  "20.128.0.0", "Calculated Next Address Correctly")	
}

func TestGetUserInput(t *testing.T) {
	assert := assert.New(t)
	Address := os.Getenv("NETWORK")
        os.Setenv("NETWORK", "")
        Vnid := os.Getenv("VNID")
        os.Setenv("VNID", "")
        Vtep := os.Getenv("REMOTE_VTEPIP")
        os.Setenv("REMOTE_VTEPIP", "")
        func() {
                defer func() {
                        if r := recover(); r == nil {
                                t.Errorf("GetUserInput should have panicked!")
                        }
                }()
                GetUserInput()
        }()
	if Address == ""{
		Address = "192.168.254.254/24"
	}
	if Vnid == ""{
		Vnid = "200"
	}
	if Vtep == ""{
		Vtep = "10.10.10.10"
	}	
        os.Setenv("NETWORK", Address)
        os.Setenv("VNID", Vnid)
        os.Setenv("REMOTE_VTEPIP", Vtep)
        input := GetUserInput()
	assert.Equal(input.Address,  Address, "Parsed User succesfully")	
	assert.Equal(input.Vnid,  Vnid, "Parsed User succesfully")	
	assert.Equal(input.RemoteVtepIP,  Vtep, "Parsed User succesfully")	
        os.Setenv("NETWORK", "300.300.300.300/300")
        func() {
                defer func() {
                        if r := recover(); r == nil {
                                t.Errorf("GetUserInput should have panicked!")
                        }
                }()
                GetUserInput()
        }()
        os.Setenv("NETWORK", Address)
}



