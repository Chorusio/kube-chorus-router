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
        "strings"
        "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (api *KubernetesAPIServer)DeleteKubeExtenderPod(originalNode v1.Node, input *Input) {
	labels := strings.Split(originalNode.Labels["NodeID"],"Node-")
	klog.Info("[INFO] node label ia", labels)
	fmt.Println("[INFO] Node labels", labels)
	if len(labels) > 0 {
		configMapName := "kube-chorus-router"
        	configMaps, err := api.Client.CoreV1().ConfigMaps(input.NameSpace).Get(configMapName, metav1.GetOptions{})
        	if err != nil {
                        return
        	}else{
			fmt.Println("[INFO] Configmap output before update", configMaps.Data)
			nodeid := "Host-"+labels[1]
			value := configMaps.Data[nodeid]
			fmt.Println("Label[0], Label[1], value", labels[0], labels[1], value)
			delete(configMaps.Data, nodeid);
			delete(configMaps.Data, "Node-"+value);
			delete(configMaps.Data, "Mac-"+value);
			delete(configMaps.Data, "Interface-"+value);
			delete(configMaps.Data, "CNI-"+value);
			fmt.Println("[INFO] Configmap output after updation", configMaps.Data)
        		configMaps, _ = api.Client.CoreV1().ConfigMaps(input.NameSpace).Update(configMaps)
		}
	}
}

func Delete(){
	klog.InitFlags(nil)
	input := GetUserInput()
	api,_ := CreateK8sApiserverClient()
	nameSpace, _  := api.CreateK8sNameSpace()
	input.NameSpace = nameSpace
	api.CreateClusterRoles(input)
	api.CreateClusterRoleBindings(input)
	serviceAccount, _  := api.CreateK8sServiceAccount(input)
	input.ServiceAccount = serviceAccount
	api.DeletePods(input)
	api.DeleteK8sConfigMap(input)
	
}

func (api *KubernetesAPIServer)DeleteK8sConfigMap(input *Input){
	configMapName := "kube-chorus-router"
	api.Client.CoreV1().ConfigMaps(input.NameSpace).Delete(configMapName, metav1.NewDeleteOptions(0))
}

func (api *KubernetesAPIServer)DeletePods(input *Input){
	klog.Infof("Cleaning stale pods created by kube-chorus-router...")
	pods, _ := api.Client.CoreV1().Pods(input.NameSpace).List(metav1.ListOptions{})
	for i := 0; i < len(pods.Items); i++ {
		if (strings.Contains(pods.Items[i].Name, "kube-chorus-router")){
        		api.Client.CoreV1().Pods(input.NameSpace).Delete(pods.Items[i].Name, metav1.NewDeleteOptions(90))
		}
	}
}
