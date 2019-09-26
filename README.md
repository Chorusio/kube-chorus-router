# kube-router
**Chorus kube router** is a solution to establish route between Kubernetes  cluster nodes and non Kubernetes nodes.  This is use-full for ingress resource where  service IP's (Pod IP) can be configured on Ingress device for load balancing front end applications. kube-router can be used for creating a route between Kubernetes and Ingress device. A service of type nodeport and loadbalncer  can reach to the pods irrespective of network and subnet. However such solution reduces the performance due to translation and load balancing at different level. So typically customer prefer ingress mechanism [Reach to pod directly from external load balancer] to expose the route to the service. **Chorus kube router** helps to create network connectivity for such kind of deployments. **kube-router** is a micro service provided by **Chorus** that helps to create network between the Kubernetes cluster nodes and non kubernetes aware devices [F5, A10, Citrix ADC]. kube-router takes care of networking changes on kubernetes and produces a configmap output which can be used by vendors to establish route from thier devices to kubernetes cluster.

[![Build Status](https://travis-ci.com/Chorusio/K8s-Route-Extender.svg?token=GfEuWKxn7TJJesWboygR&branch=master)](https://travis-ci.com/Chorusio/K8s-Route-Extender)
[![codecov](https://codecov.io/gh/Chorusio/K8s-Route-Extender/branch/master/graph/badge.svg?token=9c5R8ukQGY)](https://codecov.io/gh/Chorusio/K8s-Route-Extender)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./license/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Chorusio/K8s-Route-Extender)](https://goreportcard.com/report/github.com/Chorusio/K8s-Route-Extender)
[![Docker Repository on Quay](https://quay.io/repository/chorus/chorus-kube-router/status "Docker Repository on Quay")](https://quay.io/repository/chorus/chorus-kube-router)
[![GitHub stars](https://img.shields.io/github/stars/Chorusio/K8s-Route-Extender.svg)](https://github.com/Chorusio/K8s-Route-Extender/stargazers)
[![HitCount](http://hits.dwyl.io/Chorusio/K8s-Route-Extender.svg)](http://hits.dwyl.io/Chorusio/K8s-Route-Extender)

---


# Contents

-  [Overview](#overview)
-  [Architecture](#architecture)
-  [How it works](#how-it-works)
-  [Get started](#get-started)
-  [Issues](#issues)
-  [Code of conduct](#code-of-conduct)
-  [License](#License)

# Overview

In Kubernetes environments, when you expose the services for external access through the ingress device, to route the traffic into the cluster, you need to appropriately configure the network between the Kubernetes nodes and the Ingress device. Configuring the network is challenging as the pods use private IP addresses based on the CNI framework. Without proper network configuration, the Ingress device cannot access these private IP addresses. Also, manually configuring the network to ensure such reachability is cumbersome in Kubernetes environments. 

**Chorus** provides a microservice called as **kube-router** that you can use to create the network between the cluster and the Ingress devices.

# Architecture

The following diagram provides the high-level architecture of the kube-router:

![](./docs/images/kube-router.png)

**kube-router** creates a seperate network for any external devices and generate config-map file with network details. It does the following 
- Manage seperate subnet for non kubernetes aware nodes
- Creates  vxlan overlays for the external non kubernetes aware nodes
- Genrate a config-map file which can be used for creating other endpoint overlays
# How it works

**kube-router** creates a route entry point in each node present in the kubernetes cluster. When a node leaves it removes the route entry on the node. This information keeps in configmap which can be used for extending the route  with other nodes. Config map can be found in kube-system namespace with the endpoint details.



```
MacBook-Pro:k8s-route-extender$ kubectl get cm -n kube-system kube-chorus-router -o json
{
    "apiVersion": "v1",
    "data": {
        "CNI-10.106.170.62": "10.244.1.1/24",
        "CNI-10.106.170.63": "10.244.6.1/24",
        "EndpointIP": "192.168.1.254",
        "Host-cb716e61-cab6-437e-a84a-d26a908260bc": "10.106.170.62",
        "Host-d666ca12-5b2e-4716-a243-ece13e780122": "10.106.170.63",
        "Interface-10.106.170.62": "192.168.254.1",
        "Interface-10.106.170.63": "192.168.254.2",
        "Mac-10.106.170.62": "76:13:e1:c7:4b:f6",
        "Mac-10.106.170.63": "b2:30:00:b1:88:49",
        "Node-10.106.170.62": "10.106.170.62",
        "Node-10.106.170.63": "10.106.170.63"
    },
    "kind": "ConfigMap",
    "metadata": {
        "creationTimestamp": "2019-09-25T12:02:26Z",
        "name": "kube-chorus-router",
        "namespace": "kube-system",
        "resourceVersion": "5439136",
        "selfLink": "/api/v1/namespaces/kube-system/configmaps/kube-chorus-router",
        "uid": "20f8407b-3871-4595-9bc4-d4eb65fb80b8"
    }
}
MacBook-Pro:k8s-route-extender$ 
``` 
There are two hosts present in the given cluster, which can be identified by Host tag [```Host-cb716e61-cab6-437e-a84a-d26a908260bc``` and ```Host-d666ca12-5b2e-4716-a243-ece13e780122```] both having values as 10.106.170.62 and 10.106.170.63 respectively. These are nodes IP in the cluster. kube router creates an interface for each node which has subnet of 192.168.254.1 and 192.168.254.2 which maps to CNI subnet of 10.244.1.1 10.244.6.1 respectively. 
# Get started

Chorus kube-router can be used in the following two ways:

-  In cluster configuration. In this configuration, kube-router is run as **microservice**.
-  Out of the cluster configuration. In this configuration, the chorus is run as a **process**.

  
## Using kube-router as a process

Before you deploy the kube-router package, ensure that you have installed Go binary for running kube-router.

Perform the following:

1.  Download or clone the `kube-router` package.

2.  Navigate to the build directory 

3.  Start the `kube-router` using `make run`


## Using kube-router as a microservice

Refer the [deployment](deploy/README.md) page for running kube-router as a microservice inside the Kubernetes cluster.


# Issues

Use github issue template to report any bug. Describe the bug in details and capture the logs and share.

# Code of conduct

This project adheres to the [Kubernetes Community Code of Conduct](https://github.com/kubernetes/community/blob/master/code-of-conduct.md). By participating in this project you agree to abide by its terms.

# License

[Apache License 2.0](./license/LICENSE)
