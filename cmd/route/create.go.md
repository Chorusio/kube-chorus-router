package route

import (
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
	"path/filepath"
	"strconv"
)

var (
	kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config     *restclient.Config
	err        error
	podcount   = 0
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

func Create(){
	api,_ := CreateK8sApiserverClient()
	CreateKubeRouteExtender(api, "citrix", "1.1.1.1", "255.255.255.0", "aa:bb:cc:dd:ff:ee", "2.2.2.2")
}

func CreateKubeRouteExtender(api *KubernetesAPIServer, namespace string, network string, nexthop string, mac string, vtepip string) {
	func GenerateNodeNetworkInfo(api *KubernetesAPIServer, obj interface{}, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput, node *Node, originalNode v1.Node, PodCIDR string) {
        podcount = podcount + 1
        klog.Info("[INFO] Generating PODCIDR and Node Information")
        patchBytes := []byte(fmt.Sprintf(`{"metadata":{"labels":{"NodeIP":"%s"}}}`, node.IPAddr))
        if node.IPAddr == "" {
                patchBytes = []byte(fmt.Sprintf(`{"metadata":{"labels":{"NodeIP":"%s"}}}`, node.ExternalIPAddr))
        }
        time.Sleep(10 * time.Second) //TODO, We have to wait till Node is available.
        if _, err = api.Client.CoreV1().Nodes().Patch(originalNode.Name, types.StrategicMergePatchType, patchBytes); err != nil {
                klog.Errorf("[ERROR] Failed to Patch label %v", err)
        } else {
                klog.Info("[INFO] Updated node  label")
        }
        command := []string{"/bin/bash", "-c"}
        args := []string{
                "vtepmac=`ifconfig flannel.1 | grep -o -E '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}' `; echo \"InterfaceInfo ${vtepmac}\"; theIPaddress=`ip -4 addr show flannel.1  | grep inet | awk '{print $2}' | cut -d/ -f1`;  hostip=`ip -4 addr show eth0  | grep inet | awk '{print $2}' | cut -d/ -f1`; echo \"IP Addredd ${theIPaddress}\"; echo \"Host IP Address ${hostip}\"; `kubectl patch configmap citrix-node-controller  -p '{\"data\":{\"'\"$theIPaddress\"'\": \"'\"$vtepmac\"'\"}}'`;  `kubectl patch configmap citrix-node-controller  -p '{\"data\":{\"'\"$hostip\"'\": \"'\"$theIPaddress\"'\"}}'`;  ip route add ${network}  via  ${nexthop} dev flannel.1 onlink; arp -s ${nexthop}  ${ingmac}  dev flannel.1; bridge fdb add ${ingmac} dev flannel.1 dst ${vtepip}; sleep 3d;"}

	SecurityContext := new(v1.SecurityContext)
        Capabilities := new(v1.Capabilities)
        Capabilities.Add = append(Capabilities.Add, "NET_ADMIN")
        SecurityContext.Capabilities = Capabilities
        pod := &v1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                        Name: "citrixdummypod" + strconv.Itoa(podcount),
                },
                Spec: v1.PodSpec{
                        ServiceAccountName: "citrix-node-controller",
                        HostNetwork:        true,
                        Containers: []v1.Container{
                                {
                                        Name:            "citrixdummypod" + strconv.Itoa(podcount),
                                        Image:           "quay.io/citrix/dummynode:latest",
                                        Command:         command,
                                        Args:            args,
                                        SecurityContext: SecurityContext,
                                        Env: []v1.EnvVar{
                                                {Name: "network", Value: ControllerInputObj.IngressDevicePodSubnet},
                                                {Name: "nexthop", Value: ControllerInputObj.IngressDevicePodIP},
                                                {Name: "ingmac", Value: ControllerInputObj.IngressDeviceVtepMAC},
                                                {Name: "vtepip", Value: ControllerInputObj.IngressDeviceVtepIP},
                                        },
                                },
                        },
                },
        }
        nodeSelector := make(map[string]string)
        nodeSelector["NodeIP"] = node.IPAddr
        pod.Spec.NodeSelector = nodeSelector
        //time.Sleep(10 * time.Second) //TODO, We have to wait till Pod is available.
        //if _, err = api.Client.CoreV1().Pods(ControllerInputObj.Namespace).Create(pod); err != nil {
        //      klog.Error("Failed to Create a Pod " + err.Error())
        //}
        time.Sleep(60 * time.Second) //TODO, We have to wait till Node is available.

        //pod, err = api.Client.CoreV1().Pods(ControllerInputObj.Namespace).Get(pod.Name, metav1.GetOptions{})
        //if err != nil {
	{
        //      fmt.Errorf("pod Get API error: %v \n pod: %v", err, pod)
        //}
        configMaps, err := api.Client.CoreV1().ConfigMaps(ControllerInputObj.Namespace).Get("citrix-node-controller", metav1.GetOptions{})
        if err != nil {
	{
                fmt.Errorf("ConfigMap Get API error")
        }
        if configMaps != nil {
                ConfigMapData := make(map[string]string)
                ConfigMapData = configMaps.Data
                klog.Info("CONFIG MAP DATA", ConfigMapData)
                node.PodAddress = ConfigMapData[node.IPAddr]
                node.PodVTEP = ConfigMapData[node.PodAddress]
                node.PodNetMask = ConvertPrefixLenToMask("24")
                node.PodMaskLen = "24"
        } else {
                klog.Error("Config MAP is Empty \n")
        }





	command := []string{"/bin/bash", "-c"}
	args := []string{
		"vtepmac=`ifconfig flannel.1 | grep -o -E '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}' `; echo \"InterfaceInfo ${vtepmac}\"; theIPaddress=`ip -4 addr show flannel.1  | grep inet | awk '{print $2}' | cut -d/ -f1`;  hostip=`ip -4 addr show eth0  | grep inet | awk '{print $2}' | cut -d/ -f1`; echo \"IP Addredd ${theIPaddress}\"; echo \"Host IP Address ${hostip}\"; `kubectl patch configmap citrix-node-controller  -p '{\"data\":{\"'\"$theIPaddress\"'\": \"'\"$vtepmac\"'\"}}'`;  `kubectl patch configmap citrix-node-controller  -p '{\"data\":{\"'\"$hostip\"'\": \"'\"$theIPaddress\"'\"}}'`;  ip route add ${network}  via  ${nexthop} dev flannel.1 onlink; arp -s ${nexthop}  ${ingmac}  dev flannel.1;bridge fdb add ${ingmac} dev flannel.1 dst ${vtepip}; sleep 3d;"}

	SecurityContext := new(v1.SecurityContext)
	Capabilities := new(v1.Capabilities)
	Capabilities.Add = append(Capabilities.Add, "NET_ADMIN")
	SecurityContext.Capabilities = Capabilities

	DaemonSet := &appv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "routeaddpod",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "routeaddpod",
			},
		},
		Spec: appv1.DaemonSetSpec{
			MinReadySeconds: 2,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "routeaddpod",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "routeaddpod",
					},
					Name: "routeaddpod",
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "citrix-node-controller",
					HostNetwork:        true,
					Containers: []v1.Container{
						{
							Name:            "citrixdummypod" + strconv.Itoa(podcount),
							Image:           "quay.io/citrix/dummynode:latest",
							Command:         command,
							Args:            args,
							SecurityContext: SecurityContext,
							Env: []v1.EnvVar{
								{Name: "network", Value: network},
								{Name: "nexthop", Value: nexthop},
								{Name: "ingmac", Value: mac},
								{Name: "vtepip", Value: vtepip},
							},
						},
					},
				},
			},
		},
	}
	_, err := api.Client.AppsV1().DaemonSets(namespace).Create(DaemonSet)
	if err != nil {
		klog.Error("[ERROR] Failed to create daemon set:", err)
	}
}
