package route
/*
func DeleteKubeRouteExtender(api *KubernetesAPIServer) {
	command := []string{"/bin/bash", "-c"}
	args := []string{
		"vtepmac=`ifconfig flannel.1 | grep -o -E '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}' `; echo \"InterfaceInfo ${vtepmac}\"; theIPaddress=`ip -4 addr show flannel.1  | grep inet | awk '{print $2}' | cut -d/ -f1`;  hostip=`ip -4 addr show eth0  | grep inet | awk '{print $2}' | cut -d/ -f1`; echo \"IP Addredd ${theIPaddress}\"; echo \"Host IP Address ${hostip}\";ip route delete ${network}  via  ${nexthop} dev flannel.1 onlink; arp -d ${nexthop}  dev flannel.1; bridge fdb delete ${ingmac} dev flannel.1 dst ${vtepip}; sleep 3d;"}

	SecurityContext := new(v1.SecurityContext)
	Capabilities := new(v1.Capabilities)
	Capabilities.Add = append(Capabilities.Add, "NET_ADMIN")
	SecurityContext.Capabilities = Capabilities

	DaemonSet := &appv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "citrixroutecleanuppod",
			Namespace: ControllerInputObj.Namespace,
			Labels: map[string]string{
				"app": "citrixroutecleanuppod",
			},
		},
		Spec: appv1.DaemonSetSpec{
			MinReadySeconds: 2,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "citrixroutecleanuppod",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "citrixroutecleanuppod",
					},
					Name: "citrixroutecleanuppod",
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "citrix-node-controller",
					HostNetwork:        true,
					Containers: []v1.Container{
						{
							Name:            "citrixdummypod" + strconv.Itoa(podcount),
							Image:           "quay.io/citrix/dummynode:latest",
							Command:         command,
							Args:            args,
							SecurityContext: SecurityContext,
							Env: []v1.EnvVar{
								{Name: "network", Value: ControllerInputObj.IngressDevicePodSubnet},
								{Name: "nexthop", Value: ControllerInputObj.IngressDevicePodIP},
								{Name: "ingmac", Value: ControllerInputObj.IngressDeviceVtepMAC},
								{Name: "vtepip", Value: ControllerInputObj.IngressDeviceVtepIP},
							},
						},
					},
*/
