// Copyright 2019 Chorus  authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package route

import (
        "fmt"
	"k8s.io/klog"	
	"os"
	"strings"
	"strconv"
	"encoding/binary"
)

type Input struct{
	Mode string
	NameSpace string
	ServiceAccount string
	ConfigMap string
	Address string
	NodeIP string
	Network string
	PrefixLen string
	Netmask string
	Vnid	string
	RemoteVtepIP string
	RemoteTunnelEndPoint string
	TunnelNetwork string
	Type string
	RemoteIP string
	NextAddress string
	NodeCIDR string
	VxlanPort string
}

func ValidateAddress(address string)bool{
	ipaddress := strings.Split(address, ".")
        firstOctect, err := strconv.Atoi(ipaddress[0])
        if err != nil {
                return false
        }
        if firstOctect < 0 || firstOctect > 255 {
                return false
        }
        secondOctect, err := strconv.Atoi(ipaddress[1])
        if err != nil {
                return false
        }
        if secondOctect < 0 || secondOctect > 255 {
                return false
        }
        thirdOctect, err := strconv.Atoi(ipaddress[2])
        if err != nil {
                return false
        }
        if thirdOctect < 0 || thirdOctect > 255 {
                return false
        }
        fourthOctect, err := strconv.Atoi(ipaddress[3])

        if err != nil {
                return false
        }
        if fourthOctect < 0 || fourthOctect > 255 {
                return false
        }
        return true
}

// ConvertPrefixLenToMask convert the prefix len to netmask (dotted) format.
func ConvertPrefixLenToMask(prefixLen string) string {
        len, _ := strconv.Atoi(prefixLen)
        netmask := (uint32)(^(1<<(32-(uint32)(len)) - 1))
        bytes := make([]byte, 4)
        binary.BigEndian.PutUint32(bytes, netmask)
        netmaskdot := fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3])
        return netmaskdot
}


func ValidatePrefixLen(prefixlen string)bool{
	len,_ := strconv.Atoi(prefixlen)
	if (len <=0 || len>=32){
		return false
	}
	return true
}

func ExtractNetworkAndPrefix(address string)(string, string){
	data := strings.Split(address, "/") 
	prefixLen := data[1]
	network   := data[0]
	return network, prefixLen
}

var PrefixSubnetTable = make(map[string]string)
func InitPrefixSubnetTable(){
        PrefixSubnetTable["0"] = "0"
        PrefixSubnetTable["1"] = "128"
        PrefixSubnetTable["2"] = "192"
        PrefixSubnetTable["3"] = "224"
        PrefixSubnetTable["4"] = "240"
        PrefixSubnetTable["5"] = "248"
        PrefixSubnetTable["6"] = "252"
        PrefixSubnetTable["7"] = "254"
        PrefixSubnetTable["8"] = "0"
        PrefixSubnetTable["9"] = "128"
        PrefixSubnetTable["10"] = "192"
        PrefixSubnetTable["11"] = "224"
        PrefixSubnetTable["12"] = "240"
        PrefixSubnetTable["13"] = "248"
        PrefixSubnetTable["14"] = "252"
        PrefixSubnetTable["15"] = "254"
        PrefixSubnetTable["16"] = "0"
        PrefixSubnetTable["17"] = "128"
        PrefixSubnetTable["18"] = "192"
        PrefixSubnetTable["19"] = "224"
        PrefixSubnetTable["20"] = "240"
        PrefixSubnetTable["21"] = "248"
        PrefixSubnetTable["22"] = "252"
        PrefixSubnetTable["23"] = "254"
        PrefixSubnetTable["24"] = "0"
        PrefixSubnetTable["25"] = "128"
        PrefixSubnetTable["26"] = "192"
        PrefixSubnetTable["27"] = "224"
        PrefixSubnetTable["28"] = "240"
        PrefixSubnetTable["29"] = "248"
        PrefixSubnetTable["30"] = "252"
        PrefixSubnetTable["31"] = "254"
        PrefixSubnetTable["32"] = "255"
}


func InitializeNodeIP(input *Input){
	klog.Infof("Initializing Kube Router")
	InitPrefixSubnetTable()
	SubnetMasked := strings.Split(input.Network, ".")
        var network int
        var Network string
	var HostMax string
	var RemoteIP string
	prefix,_ := strconv.Atoi(input.PrefixLen)
        if (prefix >= 24){
                tmp_start := SubnetMasked[3]
                tmp, _ := strconv.Atoi(tmp_start)
                tmp2,_ := strconv.Atoi(PrefixSubnetTable[input.PrefixLen])
                network = tmp & tmp2
                Network = SubnetMasked[0]+"."+SubnetMasked[1]+"."+SubnetMasked[2]+"."+strconv.Itoa(network)
		HostMax = SubnetMasked[0]+"."+SubnetMasked[1]+"."+SubnetMasked[2]+"."+strconv.Itoa(255-network)
		RemoteIP = SubnetMasked[0]+"."+SubnetMasked[1]+"."+SubnetMasked[2]+"."+strconv.Itoa(255-network-1)
        } else if (prefix >= 16){
                tmp_start := SubnetMasked[2]
                tmp, _ := strconv.Atoi(tmp_start)
                tmp2,_ := strconv.Atoi(PrefixSubnetTable[input.PrefixLen])
                network = tmp & tmp2
                Network = SubnetMasked[0]+"."+SubnetMasked[1]+"."+strconv.Itoa(network)+".0"
		HostMax = SubnetMasked[0]+"."+SubnetMasked[1]+"."+strconv.Itoa(255-network)+".255"
		RemoteIP = SubnetMasked[0]+"."+SubnetMasked[1]+"."+strconv.Itoa(255-network-1)+".254"
        }else if (prefix >= 8){
                tmp_start := SubnetMasked[1]
                tmp, _ := strconv.Atoi(tmp_start)
                tmp2,_ := strconv.Atoi(PrefixSubnetTable[input.PrefixLen])
                network = tmp & tmp2
                Network = SubnetMasked[0]+"."+strconv.Itoa(network)+".0.0"
		HostMax = SubnetMasked[0]+"."+strconv.Itoa(255-network)+"255.255"
		RemoteIP = SubnetMasked[0]+"."+strconv.Itoa(255-network-1)+"254.254"
        }
	klog.Infof("Host max is %v and Remote IP is %v", HostMax, RemoteIP)
	input.NextAddress = Network
	input.RemoteIP = RemoteIP
}

func GetUserInput() (*Input){
	input := new(Input)
	input.Mode = "Dev"
	if (os.Getenv("MODE") == "Test"){
		input.Mode = "Test"
	}
        configError := 0
	input.Type = os.Getenv("TYPE")
	if (len(input.Type) == 0) {
		input.Type = "VXLAN"
		klog.Infof("Type (TYPE) is not provided considering deafult as VXLAN")
	}
	if (input.Type == "VXLAN"){
		input.Address = os.Getenv("NETWORK")
        	if len(input.Address) == 0 {
                	klog.Infof("New Private Subnet (NETWORK Eg 192.168.1.0/16) is must for extending the route")
                	configError = 1
        	}else {
			input.Network, input.PrefixLen = ExtractNetworkAndPrefix(input.Address)
			if !(ValidateAddress(input.Network)){
                		klog.Infof("Invalid Address format (NETWORK Eg 192.168.1.2)")
                		configError = 1
			}
			if !(ValidatePrefixLen(input.PrefixLen)){
                		klog.Infof("Invalid Address prefix len (0>prefixlen<32)")
                		configError = 1
			}
		}
	
		input.Vnid = os.Getenv("VNID")
        	if len(input.Vnid) == 0 {
                	klog.Infof("A unique VNID (VNID) is must for extending the route")
                	configError = 1
        	}
		input.VxlanPort = os.Getenv("VXLAN_PORT")
        	if len(input.VxlanPort) == 0 {
                	klog.Infof("VxlanPort (VXLAN_PORT) is must for extending the route")
                	configError = 1
        	}
		input.RemoteVtepIP = os.Getenv("REMOTE_VTEPIP")
        	if len(input.RemoteVtepIP) == 0 {
                	klog.Infof("A unique REMOTE_VTEPIP (REMOTE_VTEPIP) is must for extending the route")
                	configError = 1
        	}
	}else if (input.Type == "IPIP"){
		input.RemoteTunnelEndPoint = os.Getenv("REMOTE_TUNNEL_ENDPOINT")
        	if len(input.RemoteTunnelEndPoint) == 0 {
			klog.Infof("A unique REMOTE_TUNNEL_ENDPOINT (REMOTE_TUNNEL_ENDPOINT) is must for extending the route")
                	configError = 1
        	}
		/*
		input.TunnelNetwork = os.Getenv("TUNNEL_SUBNET")
        	if len(input.TunnelNetwork) == 0 {
			klog.Infof("TunnelNetwork (TUNNEL_SUBNET) is must for extending the route")
                	configError = 1
        	}
		*/
	}
	if configError == 1 {
                klog.Errorf("Unable to get the above mentioned input from YAML")
                panic("[ERROR] Killing Container.........Please restart  with Valid Inputs")
        }
	InitializeNodeIP(input)
	fmt.Println("Next Address", input.NextAddress)
	return input
}
