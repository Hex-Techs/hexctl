package templates

const RunTemp = `package {{.ToLower .Kind}}

import (
	"time"

	ctrl "{{.Repo}}/pkg/controller/{{.ToLower .Kind}}"
	clientset "{{.Repo}}/pkg/generated/clientset/versioned"
	informers "{{.Repo}}/pkg/generated/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

func Run(workCount int, cfg *restclient.Config, kubeClient *kubernetes.Clientset, stopCh <-chan struct{}) error {
	{{.ToLower .Kind}}Client, err := clientset.NewForConfig(cfg)
	if err != nil {
		return err
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	{{.ToLower .Kind}}InformerFactory := informers.NewSharedInformerFactory({{.ToLower .Kind}}Client, time.Second*30)

	controller := ctrl.NewController(kubeClient, {{.ToLower .Kind}}Client,
		{{.ToLower .Kind}}InformerFactory.{{.ToTitle .Group}}().{{.ToTitle .Version}}().{{.Kind}}s())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	{{.ToLower .Kind}}InformerFactory.Start(stopCh)

	if err = controller.Run(workCount, stopCh); err != nil {
		return err
	}
	return nil
}
`
