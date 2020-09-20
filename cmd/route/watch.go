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
	"encoding/json"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/api/core/v1"
	"k8s.io/klog"
        "k8s.io/client-go/tools/cache"
)

type Node struct {
	Role string
	Label string
	IPAddr string
	HostName string
	ExternalIPAddr string
	PodAddress string
	PodVTEP string
	PodNetMask string
	PodMaskLen string
}

func ParseNodeRoles(node *Node, originalNode v1.Node) {
        for _, Role := range originalNode.Spec.Taints {
                if Role.Key == "node-role.kubernetes.io/master" {
                        node.Role = "Master"
                }
        }
}

func WatchNodeEvents(api *KubernetesAPIServer, input *Input){

        nodeListWatcher := cache.NewListWatchFromClient(api.Client.CoreV1().RESTClient(), "nodes", v1.NamespaceAll, fields.Everything())
        _, nodecontroller := cache.NewInformer(nodeListWatcher, &v1.Node{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
                        CoreAddHandler(api, obj, input)
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
                        CoreUpdateHandler(api, obj, newobj)
                },
                DeleteFunc: func(obj interface{}) {
                        CoreDeleteHandler(api, obj, input)
                },
        },
        )
        stop := make(chan struct{})
        go nodecontroller.Run(stop)
        return
}

func CoreAddHandler(api *KubernetesAPIServer, obj interface{}, input *Input) {
        node, originalNode := ParseNodeEvents(api, obj, input)
        if node.Role != "Master" {
                api.CreateKubeExtenderPod(obj, node, originalNode, input)
        }
}

// ParseNodeEvents Parses the node object and store the fields to Node.
func ParseNodeEvents(api *KubernetesAPIServer, obj interface{}, input *Input) (*Node, v1.Node) {
        node := new(Node)
        node.Role = ""
        node.Label = ""
        originalObjJS, err := json.Marshal(obj)
        if err != nil {
                klog.Errorf("Failed to Marshal original object: %v", err)
        }
        var originalNode v1.Node
        if err = json.Unmarshal(originalObjJS, &originalNode); err != nil {
                klog.Errorf("Failed to unmarshal original object: %v", err)
        }
        if originalNode.Spec.Taints != nil {
                klog.Infof("Taint Information is present %v", originalNode.Spec.Taints)
                ParseNodeRoles(node, originalNode)
                klog.Infof("Setting Node Role as %v", node.Role)
        }
        return node, originalNode
}


func CoreUpdateHandler(api *KubernetesAPIServer, obj interface{}, newobj interface{}) {
	return
}
func CoreDeleteHandler(api *KubernetesAPIServer, obj interface{}, input *Input) {
        node, originalNode := ParseNodeEvents(api, obj, input)
        if node.Role != "Master" {
                api.DeleteKubeExtenderPod(originalNode, input)
        }
	return
}
