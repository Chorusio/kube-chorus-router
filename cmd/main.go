// Package main Kube Route Extender.
package main
import (
	route "k8s-route-extender/cmd/route"
)

func main(){
	route.Create()
//	route.Delete()
}
