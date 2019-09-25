package route
  
import (
	//"os"
        "testing"
	"fmt"
//	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

func (api *KubernetesAPIServer)CreateNode(){
        fmt.Println("[INFO] Creating a Dummy Node")
        NewNode := &v1.Node{
                ObjectMeta: metav1.ObjectMeta{
                        Name: "dummy",
                },
        }
	api.Client.CoreV1().Nodes().Create(NewNode)
}

func (api *KubernetesAPIServer)DeleteNode(){
        fmt.Println("[INFO] Deleting a Dummy Node")
	api.Client.CoreV1().Nodes().Delete("dummy", metav1.NewDeleteOptions(0))
}
func TestWatchNodeEvents(t *testing.T){
//	assert := assert.New(t)
	input := new(Input)
	api := CreateK8sFakeClient()
	input.NameSpace, _ = api.CreateK8sNameSpace()
	api.CreateClusterRoles(input)
	api.CreateClusterRoleBindings(input)
	api.CreateNode()
	api.DeleteNode()
}

