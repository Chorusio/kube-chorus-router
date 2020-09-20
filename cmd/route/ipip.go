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

func CreateIPIPConfigForHost()[]string{
        args := []string{
                "ifconfigdata=`ifconfig`; echo \"Interface Details ${ifconfigdata} \"; cni_name=`if [[ ${ifconfigdata} =~ \"cni\" &&  ${ifconfigdata} =~ \"flannel\" ]]; then echo \"flannel\"; elif [[ ${ifconfigdata} =~ \"cali\" &&  ${ifconfigdata} =~ \"flannel\" ]]; then echo \"cannel\"; elif [[ ${ifconfigdata} =~ \"cali\" &&  ${ifconfigdata} =~ \"tun\" ]]; then echo \"calico\"; elif [[ ${ifconfigdata} =~ \"cilium_host\" ]]; then echo \"cilium\"; else echo \"undefined\"; fi`; echo \"CNI Name is ${cni_name} \"; ethName=`ip route | grep  default | awk '$4 == \"dev\"{ print $5 }'`; ip link delete routertun0; ifNameA=`ifconfig | grep cni | head -n 1 | awk '{print $1}' | sed 's/://g'`; ifNameB=`ifconfig | grep tun | head -n 1 | awk '{print $1}' | sed 's/://g'`; ifNameC=`if [[ ${cni_name} =~ \"cannel\" ]]; then ifconfig | grep flannel | head -n 1 | awk '{print $1}' | sed 's/://g'`;  ifNameD=`if [[ ${cni_name} =~ \"cilium\" ]]; then ifconfig | grep cilium_host | head -n 1 | awk '{print $1}' | sed 's/://g'; else echo \"unknown\"; fi`; echo \"IfnameA ${ifNameA} IfNameB ${ifNameB} IfNameC ${ifNameC} IfNameD ${ifNameD}\"; ifName=`if [[ ${ifNameA} =~ \"cni\" ]]; then echo ${ifNameA}; elif  [[ ${ifNameB} =~ \"tun\" ]]; then  echo ${ifNameB}; elif  [[ ${ifNameC} =~ \"flannel\" ]]; then echo ${ifNameC};  elif  [[ ${ifNameD} =~ \"cilium_host\" ]]; then echo ${ifNameD}; else echo \"undefined\"; fi`; echo \"COnfigmap Name is ${configMap} and Namespace is ${nameSpace}\"; echo \"Host Interface ${ethName}\"; echo \"CNI Interface ${ifName}\"; hostip=`ip -4 addr show $ethName  | grep inet | awk '{print $2}' | cut -d/ -f1`; echo \"IP Addredd ${theIPaddress}\"; echo \"Host IP Address ${hostip}\";  ip tunnel add routertun0 mode ipip  remote ${remotetunnelip}  local ${hostip}; ip link set routertun0 up; `if [[ -z ${tunnelnet} ]]; then echo no network for tun ${tunnelnet}; else ip addr add ${tunnelnet} dev routertun0; fi`; `if [[ -z ${tunnelnet} ]]; then sysctl -w net.ipv4.conf.routertun0.rp_filter=0; sysctl -w net.ipv4.conf.all.rp_filter=0; sysctl -w net.ipv4.conf.default.rp_filter=0; fi`; `kubectl patch configmap ${configMap} -n ${nameSpace}  -p '{\"data\":{\"'\"Host-$nodeid\"'\": \"'\"$hostip\"'\"}}'`; `kubectl patch configmap ${configMap} -n ${nameSpace} -p '{\"data\":{\"'\"Node-$hostip\"'\": \"'\"$hostip\"'\"}}'`; sleep 3d"}

	return args
}
