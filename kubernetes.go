package netkat

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type (
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
		PodStatus      string `json:"status,omitempty"`
	}

	ServicePort struct {
		Type           string `json:"type,omitempty"`
		ClusterIP      net.IP `json:"clusterIP,omitempty"`
		ServiceName    string `json:"name,omitempty"`
		Namespace      string `json:"namespace,omitempty"`
		ExternalIP     net.IP
		AppSelector    string
		Host           string
		SourcePortName string `json:"name,omitempty"`
		Protocol       string `json:"protocol,omitempty"`
		SourcePort     int32  `json:"port,omitempty"`
		NodePort       int32  `json:"nodePort,omitempty"`
		TargetPort     int32  `json:"targetPort,omitempty"`
		TargetPortName string `json:"targetPort,omitempty"`
		IngressPath    IngressPath
		PodPort        []*PodPort
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

	KubernetesComponents struct {
		IngressPaths []*IngressPath
		ServicePorts []*ServicePort
		PodPorts     []*PodPort
	}
)

func (co *KubernetesComponents) FindIngressPathForHost(t *Target) (ingressPath *IngressPath, err error) {
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

func (co *KubernetesComponents) FindServicePortForHost(t *Target) (servicePort *ServicePort, err error) {
	var servicePorts []*ServicePort
	for _, s := range co.ServicePorts {
		if t.Host == s.Host && t.Port == s.SourcePort && t.IpAddress.Equal(s.ExternalIP) {
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

func (co *KubernetesComponents) FindServicePortForIngressPath(i *IngressPath) (servicePort *ServicePort, err error) {
	var servicePorts []*ServicePort
	for _, s := range co.ServicePorts {
		if i.Namespace == s.Namespace && i.ServiceName == s.ServiceName && (i.ServiceIntPort == s.SourcePort || i.ServiceStrPort == s.SourcePortName) {
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

func (co *KubernetesComponents) FindPodPortForServicePort(s *ServicePort) (podPorts []*PodPort, err error) {
	for _, p := range co.PodPorts {
		if s.Namespace == p.Namespace && s.AppSelector == p.App && (s.TargetPort == p.ContainerPort || s.TargetPortName == p.PortName) {
			podPorts = append(podPorts, p)
		}
	}
	if len(podPorts) == 0 {
		err = errors.New("could not find pod port matching the service port")
	}
	return
}

func (co *KubernetesComponents) FindServicePortForPodPort(p *PodPort) (servicePort *ServicePort, err error) {
	var servicePorts []*ServicePort
	for _, s := range co.ServicePorts {
		if p.Namespace == s.Namespace && p.App == s.AppSelector && (p.ContainerPort == s.TargetPort || p.PortName == s.TargetPortName) {
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

func (co *KubernetesComponents) FindIngressPathForServicePort(s *ServicePort) (ingressPath *IngressPath, err error) {
	var ingressPaths []*IngressPath
	for _, i := range co.IngressPaths {
		if i.Namespace == s.Namespace && i.ServiceName == s.ServiceName && (i.ServiceIntPort == s.SourcePort || i.ServiceStrPort == s.SourcePortName) {
			ingressPaths = append(ingressPaths, i)
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

func (c *Client) GetComponents() (components *KubernetesComponents) {
	pods := c.GetPods()
	svcs := c.GetServices()
	ings := c.GetIngresses()
	components = &KubernetesComponents{
		IngressesToIngressPaths(ings),
		ServicesToServicePorts(svcs),
		PodsToPodPorts(pods),
	}
	return
}

func (c *Client) GetPods() (apiPods *v1.PodList) {
	apiPods, err := c.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
	}
	return
}

func PodsToPodPorts(apiPods *v1.PodList) (podPorts []*PodPort) {
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
						PodStatus:      string(pod.Status.Phase),
					},
				)
			}
		}
	}
	return
}

func (c *Client) GetServices() (apiServices *v1.ServiceList) {
	apiServices, err := c.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
	}
	return
}

func ServicesToServicePorts(apiServices *v1.ServiceList) (servicePorts []*ServicePort) {
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
			var targetIntPort int32
			if port.TargetPort.IntVal == 0 && port.TargetPort.StrVal == "" {
				targetIntPort = port.Port
			} else {
				targetIntPort = port.TargetPort.IntVal
			}
			servicePorts = append(
				servicePorts,
				&ServicePort{
					ServiceName:    service.ObjectMeta.Name,
					AppSelector:    appSelector,
					Type:           string(service.Spec.Type),
					ClusterIP:      net.ParseIP(service.Spec.ClusterIP),
					ExternalIP:     ip,
					Host:           hostName,
					Namespace:      service.ObjectMeta.Namespace,
					Protocol:       string(port.Protocol),
					SourcePortName: port.Name,
					SourcePort:     port.Port,
					NodePort:       port.NodePort,
					TargetPort:     targetIntPort,
					TargetPortName: port.TargetPort.StrVal,
				},
			)

		}

	}
	return
}

func (c *Client) GetIngresses() (apiIngresses *v1beta1.IngressList) {
	apiIngresses, err := c.ExtensionsV1beta1().Ingresses("").List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
	}
	return
}

func IngressesToIngressPaths(apiIngresses *v1beta1.IngressList) (ingressPaths []*IngressPath) {
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

func (c *Client) GetPortforwardResponse(p *PodPort) (res *http.Response, err error) {
	roundTripper, upgrader, err := spdy.RoundTripperFor(c.Config)
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
		return
	}
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", p.Namespace, p.PodName)
	hostIP := strings.TrimLeft(c.Config.Host, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)
	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%v", p.ContainerPort)}, stopChan, readyChan, out, errOut)
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
		return
	}
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		go func() {
			for range readyChan {
			}
			if len(errOut.String()) != 0 {
				_ = level.Error(Logger).Log("msg", errOut.String())
				wait.Done()
				return
			} else if len(out.String()) != 0 {
				fmt.Println(out.String())
			}
			wait.Done()
		}()
		if err = forwarder.ForwardPorts(); err != nil {
			_ = level.Error(Logger).Log("msg", err)
			wait.Done()
			return
		}
	}()
	wait.Wait()
	pfUrl := fmt.Sprintf("http://127.0.0.1:%v", p.ContainerPort)
	res, err = http.Get(pfUrl)
	close(stopChan)
	return
}
