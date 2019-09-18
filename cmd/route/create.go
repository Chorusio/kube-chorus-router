package route

import (
        "fmt"
	"time"
	"k8s.io/apimachinery/pkg/types"
	"strconv"
	"strings"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbac "k8s.io/api/rbac/v1beta1"
        apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
	"path/filepath"
)

var (
	kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config     *restclient.Config
	err        error
	podcount   = 0
)

const (
        clusterRoleKind    = "ClusterRole"
        roleKind           = "Role"
        serviceAccountKind = "ServiceAccount"
        rbacAPIGroup       = "rbac.authorization.k8s.io"
)


// This is interface for Kubernetes API Server
type KubernetesAPIServer struct {
	Suffix string
	Client kubernetes.Interface
}

type QueueUpdate struct {
	Key   string
	Force bool
}


// This creates go client.
func CreateK8sApiserverClient() (*KubernetesAPIServer, error) {
	klog.Info("[INFO] Creating API Client")
	api := &KubernetesAPIServer{}
	config, err = clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		klog.Error("[WARNING] Citrix Node Controller Runs outside cluster")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.Error("[ERROR] Did not find valid kube config info")
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Error("[ERROR] Failed to establish connection")
		klog.Fatal(err)
	}
	klog.Info("[INFO] Kubernetes Client is created")
	api.Client = client
	return api, nil
}

func (api *KubernetesAPIServer)CreateK8sNameSpace()(string, error){
	nameSpace := "kube-system"
	nsObj, err := api.Client.CoreV1().Namespaces().Get(nameSpace,  metav1.GetOptions{})
	fmt.Println("Name Space object", nsObj)
	if err != nil {
		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nameSpace}}
		nsObj, err = api.Client.CoreV1().Namespaces().Create(nsSpec)
		if err == nil {
			return "kube-system", err
		}else{
			return "Error", err
		}
	}
	return "kube-system", err
}

func (api *KubernetesAPIServer)CreateK8sServiceAccount(input *Input)(string, error){
	serviceAccountName := "kube-chorus-router"
	saObj, err := api.Client.CoreV1().ServiceAccounts(input.NameSpace).Get(serviceAccountName,  metav1.GetOptions{})
	fmt.Println("Name Space object", saObj)
	if err != nil {
		serviceAccount := &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: serviceAccountName}}
		_, err := api.Client.CoreV1().ServiceAccounts(input.NameSpace).Create(serviceAccount)
		if err == nil {
			return "kube-chorus-router", err
		}else{
			return "Error", err
		}
	}
	return "kube-chorus-router", err
}

func (api *KubernetesAPIServer)CreateK8sConfigMap(input *Input)(string, error){
	configMapName := "kube-chorus-router"
	configMaps, err := api.Client.CoreV1().ConfigMaps(input.NameSpace).Get(configMapName, metav1.GetOptions{})
        if err != nil {
		configMap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: configMapName}, Data: map[string]string{"EndpointIP": input.RemoteIP},}
		configMaps, err = api.Client.CoreV1().ConfigMaps(input.NameSpace).Create(configMap)
		if err == nil {
			return "kube-chorus-router", err
		}else{
			return "Error", err
		}
	}
	fmt.Println("Config map object", configMaps)
	return "kube-chorus-router", err
}

func (api *KubernetesAPIServer)CreateClusterRoles(input *Input) error {
	Verbs := []string{"get", "list", "watch", "create", "patch"}
	Apigroups := []string{"*"}
	ApigroupsSecond := []string{""}
	ApigroupsThird := []string{"extensions"}
	Resources := []string{"configmaps"}
	clusterRole := rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kube-chorus-router",
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups:      Apigroups,
				Resources:      Resources,
				Verbs: 		Verbs,
			},
			{
				APIGroups:      ApigroupsSecond,
				Resources:      Resources,
				Verbs: 		Verbs,
			},
			{
				APIGroups:      ApigroupsThird,
				Resources:      Resources,
				Verbs: 		Verbs,
			},
		},
	}

	if _, err := api.Client.RbacV1beta1().ClusterRoles().Create(&clusterRole); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("unable to create RBAC clusterrolebinding: %v", err)
		}

		if _, err := api.Client.RbacV1beta1().ClusterRoles().Update(&clusterRole); err != nil {
			return fmt.Errorf("unable to update RBAC clusterrolebinding: %v", err)
		}
	}
	return nil
}

func (api *KubernetesAPIServer)CreateClusterRoleBindings(input *Input) error {
	clusterRoleBinding := rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kube-chorus-router",
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbacAPIGroup,
			Kind:     clusterRoleKind,
			Name:     "kube-chorus-router",
		},
		Subjects: []rbac.Subject{
			{
				Kind:      serviceAccountKind,
				Name:      "kube-chorus-router",
				Namespace: input.NameSpace,
			},
		},
	}

	if _, err := api.Client.RbacV1beta1().ClusterRoleBindings().Create(&clusterRoleBinding); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("unable to create RBAC clusterrolebinding: %v", err)
		}

		if _, err := api.Client.RbacV1beta1().ClusterRoleBindings().Update(&clusterRoleBinding); err != nil {
			return fmt.Errorf("unable to update RBAC clusterrolebinding: %v", err)
		}
	}
	return nil
}


func Create(){
	input := GetUserInput()
	api,_ := CreateK8sApiserverClient()
	nameSpace, _  := api.CreateK8sNameSpace()
	input.NameSpace = nameSpace
	api.CreateClusterRoles(input)
	api.CreateClusterRoleBindings(input)
	serviceAccount, _  := api.CreateK8sServiceAccount(input)
	input.ServiceAccount = serviceAccount
	configMap, _  := api.CreateK8sConfigMap(input)
	input.ConfigMap = configMap
	fmt.Println("Namespace", nameSpace, serviceAccount, configMap)
	WatchNodeEvents(api, input)
	if (input.Mode == "Dev"){
		select{}
	}
}


func GenerateNextAddress(input *Input)string{

	ipaddress := strings.Split(input.NextAddress, ".")
        firstOctect, _ := strconv.Atoi(ipaddress[3])
        secondOctect, _ := strconv.Atoi(ipaddress[2])
        thirdOctect, _ := strconv.Atoi(ipaddress[1])
        prefix, _ := strconv.Atoi(input.PrefixLen)
	if (prefix >= 24){
		if (firstOctect<254){
			firstOctect = firstOctect+1
			input.NextAddress = ipaddress[0]+"."+ipaddress[1]+"."+ipaddress[2]+"."+strconv.Itoa(firstOctect)
		}else{
			fmt.Println("[ERROR] No More IP is  avialable in the given subnet")
		}	
	}else if (prefix >= 16){
		if (firstOctect<254){
			firstOctect = firstOctect+1
			input.NextAddress = ipaddress[0]+"."+ipaddress[1]+"."+ipaddress[2]+"."+strconv.Itoa(firstOctect)
		}else if (secondOctect < 254){
			secondOctect = secondOctect+1
			input.NextAddress = ipaddress[0]+"."+ipaddress[1]+"."+strconv.Itoa(secondOctect)+"."+ipaddress[3]
		}else{
			fmt.Println("[ERROR] No More IP is  avialable in the given subnet")
		}
	}else if (prefix >= 8){
		if (firstOctect<254){
			firstOctect = firstOctect+1
			input.NextAddress = ipaddress[0]+"."+ipaddress[1]+"."+ipaddress[2]+"."+strconv.Itoa(firstOctect)
		}else if (secondOctect < 254){
			secondOctect = secondOctect+1
			input.NextAddress = ipaddress[0]+"."+ipaddress[1]+"."+strconv.Itoa(secondOctect)+"."+ipaddress[3]
		}else if (thirdOctect < 254){
			thirdOctect = thirdOctect+1
			input.NextAddress = ipaddress[0]+"."+strconv.Itoa(thirdOctect)+"."+ipaddress[2]+"."+ipaddress[3]
		}else{
			fmt.Println("[ERROR] No More IP is  avialable in the given subnet")
		}
	}	
	return (input.NextAddress+"/"+input.PrefixLen)
}

func (api *KubernetesAPIServer)CreateKubeExtenderPod(obj interface{}, node *Node, originalNode v1.Node, input *Input) {
        podcount = podcount + 1
	ifip := GenerateNextAddress(input)
        klog.Info("[INFO] Generating PODCIDR and Node Information")
        patchBytes := []byte(fmt.Sprintf(`{"metadata":{"labels":{"NodeID":"%s"}}}`, "Node-"+strconv.Itoa(podcount)))
        time.Sleep(10 * time.Second) //TODO, We have to wait till Node is available.
        if _, err = api.Client.CoreV1().Nodes().Patch(originalNode.Name, types.StrategicMergePatchType, patchBytes); err != nil {
                klog.Errorf("[ERROR] Failed to Patch label %v", err)
        } else {
                klog.Info("[INFO] Updated node  label")
        }
        command := []string{"/bin/bash", "-c"}
        args := []string{
                "ifconfigdata=`ip addr show`; echo \"Interface Details ${ifconfigdata} \"; ethName=`ip route | grep  default | awk '$4 == \"dev\"{ print $5 }'`; ip link delete routervxlan0; ifNameA=`ip route | grep cni  | awk '$2 == \"dev\"{print $3}'`; ifNameB=`ip route | grep tun | grep link  | awk '$2 == \"dev\"{print $3}'`; echo \"IfnameA ${ifNameA} IfNameB ${ifNameB}\"; ifName=`if [[ ${ifNameA} =~ \"cni\" ]]; then echo ${ifNameA}; else echo ${ifNameB}; fi`; echo \"Host Interface ${ethName}\"; echo \"CNI Interface ${ifName}\";ip link add routervxlan0 type vxlan id ${vni}  dev $ethName  dstport ${vxlanPort}; ip addr add ${address} dev routervxlan0; ip link set up dev routervxlan0 mtu 1450; vtepmac=`ifconfig routervxlan0 | grep -o -E '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}' `; echo \"InterfaceInfo ${vtepmac}\"; theIPaddress=`ip -4 addr show routervxlan0  | grep inet | awk '{print $2}' | cut -d/ -f1`;  hostip=`ip -4 addr show $ethName  | grep inet | awk '{print $2}' | cut -d/ -f1`; echo \"IP Addredd ${theIPaddress}\"; echo \"Host IP Address ${hostip}\"; cniaddr=`ip -4 addr show ${ifName} | grep inet | awk '{print $2}'`; echo \"CNI IP Address ${cniaddr}\"; `kubectl patch configmap ${configMap} -n ${nameSpace}  -p '{\"data\":{\"'\"Host-$nodeid\"'\": \"'\"$hostip\"'\"}}'`;  `kubectl patch configmap ${configMap} -n ${nameSpace}  -p '{\"data\":{\"'\"Mac-$hostip\"'\": \"'\"$vtepmac\"'\"}}'`;  `kubectl patch configmap ${configMap} -n ${nameSpace} -p '{\"data\":{\"'\"Interface-$hostip\"'\": \"'\"$theIPaddress\"'\"}}'`; `kubectl patch configmap ${configMap} -n ${nameSpace} -p '{\"data\":{\"'\"CNI-$hostip\"'\": \"'\"$cniaddr\"'\"}}'`; `kubectl patch configmap ${configMap} -n ${nameSpace} -p '{\"data\":{\"'\"Node-$hostip\"'\": \"'\"$hostip\"'\"}}'`; ip route add ${network}  via  ${nexthop} dev routervxlan0 onlink; bridge fdb add ${ingmac} dev routervxlan0 dst ${vtepip}; sleep 3d; iptables -A  OPENSHIFT-FIREWALL-ALLOW -p udp -m udp --dport ${vxlanPort} -m comment --comment \"VXLAN incoming\" -j ACCEPT"}


        SecurityContext := new(v1.SecurityContext)
        Capabilities := new(v1.Capabilities)
        Capabilities.Add = append(Capabilities.Add, "NET_ADMIN")
        SecurityContext.Capabilities = Capabilities
	pod := &v1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                        Name: "kube-chorus-router-" + strconv.Itoa(podcount),
                },
                Spec: v1.PodSpec{
                        ServiceAccountName: input.ServiceAccount,
                        HostNetwork:        true,
                        Containers: []v1.Container{
                                {
                                        Name:            "kube-chorus-router-" + strconv.Itoa(podcount),
                                        Image:           "quay.io/chorus/router:latest",
                                        Command:         command,
                                        Args:            args,
                                        SecurityContext: SecurityContext,
                                        Env: []v1.EnvVar{
                                                {Name: "network", Value: input.Network},
                                                {Name: "nexthop", Value: input.Netmask},
                                                {Name: "ingmac", Value: "00:00:00:00:00:00"},
                                                {Name: "vtepip", Value: input.RemoteVtepIP},
                                                {Name: "configMap", Value: input.ConfigMap},
                                                {Name: "nameSpace", Value: input.NameSpace},
                                                {Name: "vni", Value: input.Vnid},
                                                {Name: "vxlanPort", Value: input.VxlanPort},
                                                {Name: "address", Value: ifip},
                                                {Name: "nodeid", Value: strconv.Itoa(podcount)},
                                        },
                                },
                        },
                },
        }
        nodeSelector := make(map[string]string)
        nodeSelector["NodeID"] = "Node-"+strconv.Itoa(podcount)
        pod.Spec.NodeSelector = nodeSelector
	
        if _, err = api.Client.CoreV1().Pods(input.NameSpace).Create(pod); err != nil {
              klog.Error("Failed to Create a Pod " + err.Error())
        }
        time.Sleep(60 * time.Second) //TODO, We have to wait till Pod is available.

        pod, err = api.Client.CoreV1().Pods(input.NameSpace).Get(pod.Name, metav1.GetOptions{})
        if err != nil {
              fmt.Errorf("pod Get API error: %v \n pod: %v", err, pod)
        }
}
