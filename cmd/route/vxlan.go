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

func CreateVxlanConfigForHost()[]string{
        args := []string{
                "ifconfigdata=`ifconfig`; echo \"Interface Details ${ifconfigdata} \"; cni_name=`if [[ ${ifconfigdata} =~ \"cni\" &&  ${ifconfigdata} =~ \"flannel\" ]]; then echo \"flannel\"; elif [[ ${ifconfigdata} =~ \"cali\" &&  ${ifconfigdata} =~ \"flannel\" ]]; then echo \"cannel\"; elif [[ ${ifconfigdata} =~ \"cali\" &&  ${ifconfigdata} =~ \"tun\" ]]; then echo \"calico\"; elif [[ ${ifconfigdata} =~ \"cilium_host\" ]]; then echo \"cilium\"; else echo \"undefined\"; fi`; echo \"CNI Name is ${cni_name} \"; ethName=`ip route | grep  default | awk '$4 == \"dev\"{ print $5 }'`; ip link delete routervxlan0; ifNameA=`ifconfig | grep cni | head -n 1 | awk '{print $1}' | sed 's/://g'`; ifNameB=`if [[ ${cni_name} =~ \"calico\" ]]; then ifconfig | grep tun | head -n 1 | awk '{print $1}' | sed 's/://g'; fi`; ifNameC=`if [[ ${cni_name} =~ \"cannel\" ]]; then ifconfig | grep flannel | head -n 1 | awk '{print $1}' | sed 's/://g'; fi`;  ifNameD=`if [[ ${cni_name} =~ \"cilium\" ]]; then ifconfig | grep cilium_host | head -n 1 | awk '{print $1}' | sed 's/://g'; else echo \"unknown\"; fi`; echo \"IfnameA ${ifNameA} IfNameB ${ifNameB} IfNameC ${ifNameC} IfNameD ${ifNameD}\"; ifName=`if [[ ${ifNameA} =~ \"cni\" ]]; then echo ${ifNameA}; elif  [[ ${ifNameB} =~ \"tun\" ]]; then  echo ${ifNameB}; elif  [[ ${ifNameC} =~ \"flannel\" ]]; then echo ${ifNameC};  elif  [[ ${ifNameD} =~ \"cilium_host\" ]]; then echo ${ifNameD}; else echo \"undefined\"; fi`; `kubectl patch configmap ${configMap} -n ${nameSpace}  -p '{\"data\":{\"'\"VTEP\"'\": \"'\"$cni_name$ifName\"'\"}}'`; echo \"Host Interface ${ethName}\"; echo \"CNI Interface ${ifName}\";ip link add routervxlan0 type vxlan id ${vni}  dev $ethName  dstport ${vxlanPort}; ip addr add ${address} dev routervxlan0; ip link set up dev routervxlan0 mtu 1450; vtepmac=`ifconfig routervxlan0 | grep -o -E '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}' `; echo \"InterfaceInfo ${vtepmac}\"; theIPaddress=`ip -4 addr show routervxlan0  | grep inet | awk '{print $2}' | cut -d/ -f1`;  hostip=`ip -4 addr show $ethName  | grep inet | awk '{print $2}' | cut -d/ -f1`; echo \"IP Addredd ${theIPaddress}\"; echo \"Host IP Address ${hostip}\"; cniaddrA=`ip -4 addr show ${ifName} | grep inet | awk '{print $2}'`; cni_addr=`ip -4 addr show ${ifName} | grep inet | awk '{print $2}' | grep -o -E '[0-9]+[.][0-9]+[.][0-9]+[.][0-9]+'`; echo \"CNI IP Address ${cni_addr}\"; cni_pref=`ip route | grep ${ifName} | grep -o -E '/[0-9]+' | head -n 1`; echo \"CNI IP Prefix ${cni_pref}\"; cniaddr=`if [[ ${cniaddrA} =~ \"/32\" ]]; then echo ${cni_addr}${cni_pref}; else echo ${cniaddrA}; fi`; echo \"CNI Addr ${cniaddr}\";`kubectl patch configmap ${configMap} -n ${nameSpace}  -p '{\"data\":{\"'\"Host-$nodeid\"'\": \"'\"$hostip\"'\"}}'`;  `kubectl patch configmap ${configMap} -n ${nameSpace}  -p '{\"data\":{\"'\"Mac-$hostip\"'\": \"'\"$vtepmac\"'\"}}'`;  `kubectl patch configmap ${configMap} -n ${nameSpace} -p '{\"data\":{\"'\"Interface-$hostip\"'\": \"'\"$theIPaddress\"'\"}}'`; `kubectl patch configmap ${configMap} -n ${nameSpace} -p '{\"data\":{\"'\"CNI-$hostip\"'\": \"'\"$cniaddr\"'\"}}'`; `kubectl patch configmap ${configMap} -n ${nameSpace} -p '{\"data\":{\"'\"Node-$hostip\"'\": \"'\"$hostip\"'\"}}'`; ip route add ${network}  via  ${nexthop} dev routervxlan0 onlink; bridge fdb add ${ingmac} dev routervxlan0 dst ${vtepip}; iptables -D INPUT -p udp -m udp --dport ${vxlanPort} -j ACCEPT; iptables -I INPUT 1 -p udp --dport ${vxlanPort} -j ACCEPT; sleep 3d"}

	return args
}
