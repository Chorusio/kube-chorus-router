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
                klog.Errorf("[WARNING] Does not have PodCIDR Information, CNC will Generate itself")
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
                klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var originalNode v1.Node
        if err = json.Unmarshal(originalObjJS, &originalNode); err != nil {
                klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
        if originalNode.Spec.Taints != nil {
                klog.Info("[INFO] Taint Information", originalNode.Spec.Taints)
                ParseNodeRoles(node, originalNode)
                klog.Info("[INFO] Setting Node Role", node.Role)
        }
        return node, originalNode
}


func CoreUpdateHandler(api *KubernetesAPIServer, obj interface{}, newobj interface{}) {
	return
}
func CoreDeleteHandler(api *KubernetesAPIServer, obj interface{}, input *Input) {
        node, originalNode := ParseNodeEvents(api, obj, input)
        if node.Role != "Master" {
                klog.Info("[INFO] Deleting the Node from the cluster")
                api.DeleteKubeExtenderPod(obj, node, originalNode, input)
        }
	return
}
