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
	"os"
        "testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

var (
        fakeK8sApi *KubernetesAPIServer
)

func CreateK8sFakeClient()(*KubernetesAPIServer){
	if fakeK8sApi == nil {
		fake := fake.NewSimpleClientset()
		fakeK8sApi = &KubernetesAPIServer{
                	Suffix: "Test",
                	Client: fake,
        	}
	}
	return fakeK8sApi
}

func TestCreateK8sApiserverClient(t *testing.T){
	//assert := assert.New(t)
	//api, err := CreateK8sApiserverClient()
	CreateK8sApiserverClient()
        //assert.Equal(true, true, "Creating Kube Api client")
}

func TestCreateK8sNameSpace(t *testing.T){
	assert := assert.New(t)
	api := CreateK8sFakeClient()
	api.CreateK8sNameSpace()
		
	nameSpace := "kube-system"
        obj, err := api.Client.CoreV1().Namespaces().Get(nameSpace,  metav1.GetOptions{})
	fmt.Println("[TEST INFO] Namespace Data", obj)
        assert.Equal(true, (err==nil), "Name Space Creation has failed")
}

func TestCreateK8sServiceAccount(t *testing.T){
	assert := assert.New(t)
	input := new(Input)
	api := CreateK8sFakeClient()
	input.NameSpace, _ = api.CreateK8sNameSpace()
	api.CreateK8sServiceAccount(input)
	
	serviceAccountName := "kube-chorus-router"
        saObj, err := api.Client.CoreV1().ServiceAccounts(input.NameSpace).Get(serviceAccountName,  metav1.GetOptions{})
        fmt.Println("Service Account object", saObj)
        assert.Equal(true, (err==nil), "Service account Creation has failed")
}

func TestCreateK8sConfigMap(t *testing.T){
	assert := assert.New(t)
	input := new(Input)
	api := CreateK8sFakeClient()
	input.NameSpace, _ = api.CreateK8sNameSpace()
	api.CreateK8sServiceAccount(input)
	api.CreateK8sConfigMap(input)	

	configMapName := "kube-chorus-router"
        configMaps, err := api.Client.CoreV1().ConfigMaps(input.NameSpace).Get(configMapName, metav1.GetOptions{})
        fmt.Println("Config Map", configMaps)
        assert.Equal(true, (err==nil), "Config Map Creation has failed")
	_, err = api.CreateK8sConfigMap(input)	
        assert.Equal(true, (err==nil), "Config Map Creation has failed")
}

func TestCreateClusterRoles(t *testing.T){
	assert := assert.New(t)
	input := new(Input)
	api := CreateK8sFakeClient()
	input.NameSpace, _ = api.CreateK8sNameSpace()
	api.CreateClusterRoles(input)

	name := "kube-chorus-router"	
	role, err := api.Client.RbacV1beta1().ClusterRoles().Get(name, metav1.GetOptions{});

        fmt.Println("Cluster Role", role)
        assert.Equal(true, (err==nil), "Cluster ROle Creation has failed")
}

func TestCreateClusterRoleBindings(t *testing.T){
	assert := assert.New(t)
	input := new(Input)
	api := CreateK8sFakeClient()
	input.NameSpace, _ = api.CreateK8sNameSpace()
	api.CreateClusterRoles(input)
	api.CreateClusterRoleBindings(input)

	name := "kube-chorus-router"	
	role, err := api.Client.RbacV1beta1().ClusterRoleBindings().Get(name, metav1.GetOptions{});

        fmt.Println("Cluster Role Bindings", role)
        assert.Equal(true, (err==nil), "Cluster ROle bindings Creation has failed")
}

func TestCreate(t *testing.T){
	 os.Setenv("MODE", "Test")
	 Create()
}

func TestGenerateNextAddress(t *testing.T){
	assert := assert.New(t)
	input := new(Input)
	input.PrefixLen = "24"
	input.NextAddress = "20.20.20.0"
	nextAddress := GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.20.20.0 is", nextAddress)
        assert.Equal(nextAddress, "20.20.20.1/24", "Next Address Generation has failed")
	input.NextAddress = "20.20.20.254"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.20.20.254 is", nextAddress)
        assert.Equal(nextAddress, "20.20.20.254/24", "Next Address Generation has failed")
	input.PrefixLen = "16"
	input.NextAddress = "20.20.0.0"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.20.0.0 is", nextAddress)
        assert.Equal(nextAddress, "20.20.0.1/16", "Next Address Generation has failed")
	input.PrefixLen = "16"
	input.NextAddress = "20.20.0.254"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.20.0.254 is", nextAddress)
        assert.Equal(nextAddress, "20.20.1.254/16", "Next Address Generation has failed")
	input.NextAddress = "20.20.254.254"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.20.254.254 is", nextAddress)
        assert.Equal(nextAddress, "20.20.254.254/16", "Next Address Generation has failed")
	input.PrefixLen = "8"
	input.NextAddress = "20.0.0.0"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.0.0.0 is", nextAddress)
        assert.Equal(nextAddress, "20.0.0.1/8", "Next Address Generation has failed")
	input.PrefixLen = "8"
	input.NextAddress = "20.0.0.254"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.0.0.254 is", nextAddress)
        assert.Equal(nextAddress, "20.0.1.254/8", "Next Address Generation has failed")
	input.PrefixLen = "8"
	input.NextAddress = "20.0.254.254"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.0.254.254 is", nextAddress)
        assert.Equal(nextAddress, "20.1.254.254/8", "Next Address Generation has failed")
	input.NextAddress = "20.254.254.254"
	nextAddress = GenerateNextAddress(input)
	fmt.Println("[TEST INFO] Next Address generated for 20.254.254.254 is", nextAddress)
        assert.Equal(nextAddress, "20.254.254.254/8", "Next Address Generation has failed")
}

func TestCreateKubeExtenderPod(t *testing.T){
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
	var k8sNode v1.Node

	api.CreateKubeExtenderPod(nil, nil, k8sNode, input)

	//_, err = api.Client.CoreV1().Pods(input.NameSpace).Get("kube-chorus-router-1", metav1.GetOptions{})
        //assert.Equal(true, (err==nil), "Pods are created for the Nodes")
}
