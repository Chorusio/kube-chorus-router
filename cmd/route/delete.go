package route
import (
        "fmt"
	"k8s.io/klog"
        "strings"
        "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (api *KubernetesAPIServer)DeleteKubeExtenderPod(obj interface{}, node *Node, originalNode v1.Node, input *Input) {
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
