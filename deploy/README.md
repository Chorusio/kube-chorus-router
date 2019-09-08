# Deploy the kube-router

Perform the following:

1.  Download the `k8s-router.yaml` deployment file using the following command:

        wget  https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller/master/deploy/k8s-router.yaml

    The deployment file contains definitions for the following:

    -  Cluster Role (`ClusterRole`)

    -  Cluster Role Bindings (`ClusterRoleBinding`)

    -  Service Account (`ServiceAccount`)

    -  Citrix Node Controller service (`citrix-node-controller`)

    You don't have to modify the definitions for `ClusterRole`, `ClusterRoleBinding`, and `ServiceAccount` definitions. The definitions are used by Citrix node controller to monitor Kubernetes events. But, in the`citrix-node-controller` definition you have to provide the values for the environment variables that is required for Citrix k8s node controller to configure the Citric ADC.

    You must provide values for the following environment variables in the Citrix k8s node controller service definition:

    | Environment Variable | Mandatory or Optional | Description |
    | -------------------- | --------------------- | ----------- |
    | ADDRESS | Mandatory | kube-router uses this address to configure the VTEP overlay end points on nodes.| 
    | VNID | Mandatory | A unique VNID tp create a VXLAn overlays between kubernetes nodes and ingress devices.|
    | CNI_NAME | Mandatory | Provide the CNI name used in the cluster[flannel, calico, openshift-azure, etc].|
    | K8S_VXLAN_PORT | Mandatory | VXLAN PORT for overlays.|
    | REMOTE_VTEPIP | Mandatory | Ingress device VTEP IP|

1.  After you have updated the kube router deployment YAML file, deploy it using the following command:

        kubectl create -f k8s-router.yaml

1.  Apply the [config map](https://github.com/janraj/citrix-k8s-node-controller/blob/master/deploy/config_map.yaml) using the following command:

        kubectl apply -f https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller/master/deploy/config_map.yaml
