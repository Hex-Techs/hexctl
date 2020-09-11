package cluster

func StartKubernetesCluster(kc *KubernetesCluster) {
	kc.argsCheck()
	if kc.Type == "local" {
		if kc.method == "init" {
			handleWorkDir()
			renderVagrantFile(&kc.ncmd)
			startVirtualMachine()
		}
		startKubernetesCluster(kc)
	} else {
		startKubernetesCluster(kc)
	}

	if kc.method == "init" {
		kc.initCluster()
		kc.setKubeConfig()
	} else {
		kc.joinCluster()
	}
}

func startKubernetesCluster(kc *KubernetesCluster) {
	kc.swapoff()
	kc.installPackage()
	kc.setKernelConfig()
	kc.setRepo()
	kc.installClusterPackage()
	kc.setDockerConfig()
	kc.enableService()
	kc.setNodeIP()
}
