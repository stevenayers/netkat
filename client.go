package netkat

import (
	"github.com/go-kit/kit/log/level"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type (
	Client struct {
		*kubernetes.Clientset
		Config *rest.Config
	}
)

func InitClient(context string, kubeConfig string) (k8sClient Client) {
	config, err := buildConfigFromFlags(context, kubeConfig)
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
		os.Exit(1)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
		os.Exit(1)
	}
	k8sClient = Client{clientSet, config}
	return
}

func buildConfigFromFlags(context, kubeConfig string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
