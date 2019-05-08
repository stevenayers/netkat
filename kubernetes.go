package netkat

import (
	"errors"
	"github.com/go-kit/kit/log/level"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"os"
)

type (
	Client struct {
		*kubernetes.Clientset
	}

	PodPort struct {
		PodName        string `json:"name,omitempty"`
		Namespace      string `json:"namespace,omitempty"`
		App            string
		ContainerImage string `json:"image,omitempty"`
		ContainerName  string `json:"name,omitempty"`
		PortName       string `json:"name,omitempty"`
		HostPort       int32  `json:"hostPort,omitempty"`
		ContainerPort  int32  `json:"containerPort,omitempty"`
		Protocol       string `json:"protocol,omitempty"`
		HostIP         net.IP `json:"hostIP,omitempty"`
		ServicePort    ServicePort
	}

	ServicePort struct {
		Type          string `json:"type,omitempty"`
		ClusterIP     net.IP `json:"clusterIP,omitempty"`
		ServiceName   string `json:"name,omitempty"`
		Namespace     string `json:"namespace,omitempty"`
		ExternalIP    net.IP
		AppSelector   string
		Host          string
		PortName      string `json:"name,omitempty"`
		Protocol      string `json:"protocol,omitempty"`
		Port          int32  `json:"port,omitempty"`
		NodePort      int32  `json:"nodePort,omitempty"`
		TargetIntPort int32  `json:"targetPort,omitempty"`
		TargetStrPort string `json:"targetPort,omitempty"`
		IngressPath   IngressPath
		PodPort       []*PodPort
	}

	IngressPath struct {
		Host           string `json:"host,omitempty"`
		IpAddress      net.IP `json:"ipAddress,omitempty"`
		Namespace      string `json:"namespace,omitempty"`
		IngressName    string `json:"name,omitempty"`
		Path           string `json:"path,omitempty"`
		ServiceName    string `json:"serviceName,omitempty"`
		ServiceIntPort int32  `json:"servicePort,omitempty"`
		ServiceStrPort string `json:"servicePort,omitempty"`
		Service        []*ServicePort
	}

	Components struct {
		IngressPaths []*IngressPath
		ServicePorts []*ServicePort
		PodPorts     []*PodPort
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
	k8sClient = Client{clientSet}
	return
}

func buildConfigFromFlags(context, kubeConfig string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func (co *Components) FindIngressPathForHost(t *Target) (ingressPath *IngressPath, err error) {
	var ingressPaths []*IngressPath
	for _, i := range co.IngressPaths {
		if t.Host == i.Host && t.Path == i.Path && t.IpAddress.Equal(i.IpAddress) {
			ingressPaths = append(ingressPaths, i)
		}
	}
	switch {
	case len(ingressPaths) > 1:
		err = errors.New("found more than one ingress resource matching the host")
	case len(ingressPaths) == 0:
		err = errors.New("could not find ingress resource matching the host")
	default:
		ingressPath = ingressPaths[0]
	}
	return
}

func (co *Components) FindServicePortForHost(t *Target) (servicePort *ServicePort, err error) {
	var servicePorts []*ServicePort
	for _, s := range co.ServicePorts {
		if t.Host == s.Host && t.Port == s.Port && t.IpAddress.Equal(s.ExternalIP) {
			servicePorts = append(servicePorts, s)
		}
	}
	switch {
	case len(servicePorts) > 1:
		err = errors.New("found more than one service resource matching the host")
	case len(servicePorts) == 0:
		err = errors.New("could not find service resource matching the host")
	default:
		servicePort = servicePorts[0]
	}
	return
}

func (co *Components) FindServicePortForIngressPath(p *IngressPath) (servicePort *ServicePort, err error) {
	var servicePorts []*ServicePort
	for _, s := range co.ServicePorts {
		if p.Namespace == s.Namespace && p.ServiceName == s.ServiceName && (p.ServiceIntPort == s.Port || p.ServiceStrPort == s.PortName) {
			servicePorts = append(servicePorts, s)
		}
	}
	switch {
	case len(servicePorts) > 1:
		err = errors.New("found more than one service resource matching the ingress path")
	case len(servicePorts) == 0:
		err = errors.New("could not find service resource matching the ingress path")
	default:
		servicePort = servicePorts[0]
	}
	return
}

func (co *Components) FindPodPortForServicePort(s *ServicePort) (podPort *PodPort, err error) {
	var podPorts []*PodPort
	for _, p := range co.PodPorts {
		if s.Namespace == p.Namespace && s.AppSelector == p.App && (s.TargetIntPort == p.ContainerPort || s.TargetStrPort == p.PortName) {
			podPorts = append(podPorts, p)
		}
	}
	switch {
	case len(podPorts) > 1:
		err = errors.New("found more than one pod port matching the service port")
	case len(podPorts) == 0:
		err = errors.New("could not find pod port matching the service port")
	default:
		podPort = podPorts[0]
	}
	return
}

func (co *Components) FindServicePortForPodPort(p *PodPort) (servicePort *ServicePort, err error) {
	var servicePorts []*ServicePort
	for _, s := range co.ServicePorts {
		if p.Namespace == s.Namespace && p.App == s.AppSelector && (p.ContainerPort == s.TargetIntPort || p.PortName == s.TargetStrPort) {
			servicePorts = append(servicePorts, s)
		}
	}
	switch {
	case len(servicePorts) > 1:
		err = errors.New("found more than one service resource matching the pod port")
	case len(servicePorts) == 0:
		err = errors.New("could not find service resource matching the pod port")
	default:
		servicePort = servicePorts[0]
	}
	return
}

func (co *Components) FindIngressPathForServicePort(s *ServicePort) (ingressPath *IngressPath, err error) {
	var ingressPaths []*IngressPath
	for _, p := range co.IngressPaths {
		if p.Namespace == s.Namespace && p.ServiceName == s.ServiceName && (p.ServiceIntPort == s.Port || p.ServiceStrPort == s.PortName) {
			ingressPaths = append(ingressPaths, p)
		}
	}
	switch {
	case len(ingressPaths) > 1:
		err = errors.New("found more than one ingress resource matching the service")
	case len(ingressPaths) == 0:
		err = errors.New("could not find ingress resource matching the service")
	default:
		ingressPath = ingressPaths[0]
	}
	return
}

func (c *Client) GetComponents() (components *Components) {
	pods := c.GetPods()
	svcs := c.GetServices()
	ings := c.GetIngresses()
	components = &Components{ings, svcs, pods}
	return
}

func (c *Client) GetPods() (podPorts []*PodPort) {
	apiPods, err := c.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
	}
	for _, pod := range apiPods.Items {
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				appLabel, ok := pod.ObjectMeta.Labels["app"]
				if !ok {
					appLabel = ""
				}
				podPorts = append(
					podPorts,
					&PodPort{
						PortName:       port.Name,
						HostPort:       port.HostPort,
						ContainerPort:  port.ContainerPort,
						Protocol:       string(port.Protocol),
						HostIP:         net.ParseIP(port.HostIP),
						ContainerName:  container.Name,
						ContainerImage: container.Image,
						PodName:        pod.ObjectMeta.Name,
						Namespace:      pod.ObjectMeta.Namespace,
						App:            appLabel,
					},
				)
			}
		}
	}
	return
}

func (c *Client) GetServices() (servicePorts []*ServicePort) {
	apiServices, err := c.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
	}
	for _, service := range apiServices.Items {
		for _, port := range service.Spec.Ports {
			hostName, ok := service.ObjectMeta.Annotations["external-dns.alpha.kubernetes.io/hostname"]
			if !ok {
				hostName = ""
			}
			var ip net.IP
			if len(service.Status.LoadBalancer.Ingress) > 0 {
				ip = net.ParseIP(service.Status.LoadBalancer.Ingress[0].IP)
			}
			appSelector, ok := service.Spec.Selector["app"]
			if !ok {
				appSelector = ""
			}
			servicePorts = append(
				servicePorts,
				&ServicePort{
					ServiceName:   service.ObjectMeta.Name,
					AppSelector:   appSelector,
					Type:          string(service.Spec.Type),
					ClusterIP:     net.ParseIP(service.Spec.ClusterIP),
					ExternalIP:    ip,
					Host:          hostName,
					Namespace:     service.ObjectMeta.Namespace,
					PortName:      port.Name,
					Protocol:      string(port.Protocol),
					Port:          port.Port,
					NodePort:      port.NodePort,
					TargetIntPort: port.TargetPort.IntVal,
					TargetStrPort: port.TargetPort.StrVal,
				},
			)

		}

	}
	return
}

func (c *Client) GetIngresses() (ingressPaths []*IngressPath) {
	apiIngresses, err := c.ExtensionsV1beta1().Ingresses("").List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
	}
	for _, ingressResource := range apiIngresses.Items {
		for _, ingress := range ingressResource.Spec.Rules {
			for _, path := range ingress.IngressRuleValue.HTTP.Paths {
				ingressPaths = append(
					ingressPaths,
					&IngressPath{
						Path:           path.Path,
						ServiceName:    path.Backend.ServiceName,
						ServiceIntPort: path.Backend.ServicePort.IntVal,
						ServiceStrPort: path.Backend.ServicePort.StrVal,
						IngressName:    ingressResource.ObjectMeta.Name,
						IpAddress:      net.ParseIP(ingressResource.Status.LoadBalancer.Ingress[0].IP),
						Namespace:      ingressResource.ObjectMeta.Namespace,
						Host:           ingress.Host,
					},
				)
			}
		}
	}
	return
}
