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
        "testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
        "fmt"
	"k8s.io/apimachinery/pkg/types"
)
func TestDeleteKubeExtenderPod(t *testing.T){
	//assert := assert.New(t)
	input := new(Input)
	input.Address = "30.30.30.30/24"
	input.Network = "30.30.30.0"
	input.PrefixLen = "24"
	input.Vnid = "300"
	input.RemoteVtepIP = "11.11.11.11"
	input.NextAddress = "30.30.30.1"
	api := CreateK8sFakeClient()
	input.NameSpace, _ = api.CreateK8sNameSpace()
	api.CreateClusterRoles(input)
	api.CreateClusterRoleBindings(input)
	api.CreateNode()
        podcount := GenerateUUID()
        patchBytes := []byte(fmt.Sprintf(`{"metadata":{"labels":{"NodeID":"%s"}}}`, "Node-"+podcount))
        if _, err = api.Client.CoreV1().Nodes().Patch("dummy", types.StrategicMergePatchType, patchBytes); err != nil {
                klog.Errorf("[ERROR] Failed to Patch label %v", err)
        } else {
                klog.Info("[INFO] Updated node  label")
        }
	obj,_:= api.Client.CoreV1().Nodes().Get("dummy", metav1.GetOptions{})
	_, node := ParseNodeEvents(api, obj, input)
	
	configMapName := "kube-chorus-router"
        configMaps, err := api.Client.CoreV1().ConfigMaps(input.NameSpace).Get(configMapName, metav1.GetOptions{})
        if err != nil {
               klog.Info("Confug Map is empty")
        }else{
		fmt.Println("[INFO] Configmap output before update", configMaps.Data)
		configMaps.Data["EndpointIP"]="1.1.1.1";
		configMaps.Data["Host-"+podcount]="1.1.1.1";
		configMaps.Data["Node-1.1.1.1"]="1.1.1.1";
		configMaps.Data["Mac-1.1.1.1"]="aa:bb:cc:dd:ee:ff";
		configMaps.Data["Interface-1.1.1.1"]="1.1.1.1";
		configMaps.Data["CNI-1.1.1.1"]="1.1.1.1";
        	configMaps, _ = api.Client.CoreV1().ConfigMaps(input.NameSpace).Update(configMaps)
	}
	api.DeleteKubeExtenderPod(node, input)

}
