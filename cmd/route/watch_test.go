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

