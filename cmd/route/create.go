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
	"github.com/google/uuid"
        "fmt"
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
)

// GenerateUUID create a new unique number and returns in string format.
// This function uses Google UUID package to create a new UUID.
func GenerateUUID() string {
        uuid := uuid.New()
        s := uuid.String()
        return s
}

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


// CreateK8sApiserverClient creates a kubernetes client interface. 
// If the chorus router is running inside cluster it takes default Service account and generate the k8s client API interface.
// If the chorus runs outside it looks for config file in kubeconfig location and creates the client API interface.
func CreateK8sApiserverClient() (*KubernetesAPIServer, error) {
	klog.Infof("Creating Kube API Client")
	api := &KubernetesAPIServer{}
	config, err = clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		klog.Errorf("kube-chorus-router Runs outside cluster")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.Errorf("Did not find valid kube config info")
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorf("Failed to create Kube API client")
		klog.Fatal(err)
	}
	klog.Infof("Kubernetes Client is created")
	api.Client = client
	return api, nil
}

func (api *KubernetesAPIServer)CreateK8sNameSpace()(string, error){
	nameSpace := "kube-system"
	nsObj, err := api.Client.CoreV1().Namespaces().Get(nameSpace,  metav1.GetOptions{})
	if err != nil {
		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nameSpace}}
		nsObj, err = api.Client.CoreV1().Namespaces().Create(nsSpec)
		if err != nil {
			klog.Errorf("Error while creating namespace ")
			return "Error", err
		}
	}
	klog.Infof("Namespace selected is %v", nsObj.ObjectMeta.Name)
	return nsObj.ObjectMeta.Name, err
}

func (api *KubernetesAPIServer)CreateK8sServiceAccount(input *Input)(string, error){
	serviceAccountName := "kube-chorus-router"
	saObj, err := api.Client.CoreV1().ServiceAccounts(input.NameSpace).Get(serviceAccountName,  metav1.GetOptions{})
	if err != nil {
		serviceAccount := &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: serviceAccountName}}
		saObj, err = api.Client.CoreV1().ServiceAccounts(input.NameSpace).Create(serviceAccount)
		if err != nil {
			klog.Errorf("Error while creating service account ")
			return "Error", err
		}
	}
	klog.Infof("Service account selected is %v",saObj.ObjectMeta.Name)
	return saObj.ObjectMeta.Name, err
}

func (api *KubernetesAPIServer)CreateK8sConfigMap(input *Input)(string, error){
	configMapName := "kube-chorus-router"
	api.Client.CoreV1().ConfigMaps(input.NameSpace).Delete(configMapName, metav1.NewDeleteOptions(0))
	configMap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: configMapName}, Data: map[string]string{"EndpointIP": input.RemoteIP, "Type": input.Type},}
	configMaps, err := api.Client.CoreV1().ConfigMaps(input.NameSpace).Create(configMap)
	if err == nil {
		klog.Infof("Created Configmap kube-chorus-router is %v", configMaps)
		return "kube-chorus-router", err
	}else{
		klog.Errorf("Failed to create configmap kube-chorus-router %v", err)
		return "Error", err
	}
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
	klog.InitFlags(nil)
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
	var args []string 
	var ifip string 
        podcount := GenerateUUID()
	if val, ok := originalNode.ObjectMeta.Labels["NodeID"]; ok {
		tmp := strings.Split(val, "Node-")
		podcount = tmp[1]
	}
        klog.Infof("Generating PODCIDR for Node with label used is %v", "Node-"+podcount)
        patchBytes := []byte(fmt.Sprintf(`{"metadata":{"labels":{"NodeID":"%s"}}}`, "Node-"+podcount))
        if _, err = api.Client.CoreV1().Nodes().Patch(originalNode.Name, types.StrategicMergePatchType, patchBytes); err != nil {
                klog.Errorf("Failed to Patch label %v", err)
        } else {
                klog.Infof("Updated node label")
        }
	if (input.Type == "VXLAN"){
		ifip = GenerateNextAddress(input)
		args = CreateVxlanConfigForHost()
	}else if (input.Type == "IPIP"){
		args = CreateIPIPConfigForHost()
	}
        command := []string{"/bin/bash", "-c"}


        SecurityContext := new(v1.SecurityContext)
        Capabilities := new(v1.Capabilities)
        Capabilities.Add = append(Capabilities.Add, "NET_ADMIN")
        Capabilities.Add = append(Capabilities.Add, "SYS_MODULE")
        SecurityContext.Capabilities = Capabilities
        privilege := true
        SecurityContext.Privileged = &privilege
	pod := &v1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                        Name: "kube-chorus-router-" + podcount,
                },
                Spec: v1.PodSpec{
                        ServiceAccountName: input.ServiceAccount,
                        HostNetwork:        true,
                        Containers: []v1.Container{
                                {
                                        Name:            "kube-chorus-router-" + podcount,
                                        Image:           "quay.io/chorus/router:1.1.0",
                                        Command:         command,
                                        Args:            args,
                                        SecurityContext: SecurityContext,
                                        Env: []v1.EnvVar{
                                                {Name: "network", Value: input.Network},
                                                {Name: "nexthop", Value: input.Netmask},
                                                {Name: "ingmac", Value: "00:00:00:00:00:00"},
                                                {Name: "vtepip", Value: input.RemoteVtepIP},
                                                {Name: "remotetunnelip", Value: input.RemoteTunnelEndPoint},
                                                {Name: "tunnelnet", Value: input.TunnelNetwork},
                                                {Name: "configMap", Value: input.ConfigMap},
                                                {Name: "nameSpace", Value: input.NameSpace},
                                                {Name: "vni", Value: input.Vnid},
                                                {Name: "vxlanPort", Value: input.VxlanPort},
                                                {Name: "address", Value: ifip},
                                                {Name: "nodeid", Value: podcount},
                                        },
                                },
                        },
                },
        }
        nodeSelector := make(map[string]string)
        nodeSelector["NodeID"] = "Node-"+podcount
        pod.Spec.NodeSelector = nodeSelector
	
	klog.Infof("Cleaning stale pods created by kube-chorus-router...")
        err = api.Client.CoreV1().Pods(input.NameSpace).Delete(pod.Name, metav1.NewDeleteOptions(90))
	for {        
		res, err := api.Client.CoreV1().Pods(input.NameSpace).Get(pod.Name, metav1.GetOptions{})
        	if err != nil {
              		fmt.Errorf("pod Get API error: %v \n pod: %v", err, pod.Name)
			break
        	}
		if res.ObjectMeta.Name == ""{
			break
		}
	}
	klog.Infof("Creating a pod %v", pod.Name)
        if _, err = api.Client.CoreV1().Pods(input.NameSpace).Create(pod); err != nil {
              klog.Error("Failed to Create a Pod " + err.Error())
        }
	klog.Infof("Cleaning pods created by kube-chorus-router...")
        err = api.Client.CoreV1().Pods(input.NameSpace).Delete(pod.Name, metav1.NewDeleteOptions(90))
}
