# Deploy the kube-router

Perform the following:

1.  Download the `k8s-router.yaml` deployment file using the following command:

        wget  https://raw.githubusercontent.com/Chorusio/kube-chorus-router/master/deploy/k8s-router.yaml

    The deployment file contains definitions for the following:

    -  Cluster Role (`ClusterRole`)

    -  Cluster Role Bindings (`ClusterRoleBinding`)

    -  Service Account (`ServiceAccount`)

    -  Route extender service (`kube-chorus-router`)

    You don't have to modify the definitions for `ClusterRole`, `ClusterRoleBinding`, and `ServiceAccount` definitions. The definitions are used by kube-chorus-router to monitor Kubernetes events. But, in the`kube-chorus-router` definition you have to provide the values for the environment variables that is required for kube-chorus-router to configure the ADC/Router.

    You must provide values for the following environment variables in the kube-chorus-router service definition:

    **Supported Overlays**

    | Environment Variable | Mandatory or Optional | Values  |Description |
    | -------------------- | --------------------- |---------|----------- |
    | TYPE | Mandatory | VXLAN, IPIP | Type of overlay used of extending the route.| 

    **Input For VXLAN**
     
    | Environment Variable | Mandatory or Optional | Description |
    | -------------------- | --------------------- | ----------- |
    | NETWORK              | Mandatory  | kube-router uses this network to configure the VTEP overlay end points on nodes.| 
    | VNID                 | Mandatory  | A unique VNID tp create a VXLAn overlays between kubernetes nodes and ingress devices.|
    | K8S_VXLAN_PORT       | Mandatory  | VXLAN PORT for overlays.|
    | REMOTE_VTEPIP        | Mandatory  | Ingress device VTEP IP|

    **Input For IPIP**
     
    | Environment Variable | Mandatory or Optional | Description |
    | -------------------- | --------------------- | ----------- |
    | REMOTE_TUNNEL_ENDPOINT | Mandatory | Endpoint Tunnel IP.|
    | TUNNEL_SUBNET          | Optional  | Subnet for which tunnel is required. |

 
1.  After you have updated the kube router deployment YAML file, deploy it using the following command:

        kubectl create -f k8s-router.yaml

