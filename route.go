package netkat

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/goware/urlx"
	"net"
	"net/url"
	"os"
	"strconv"
)

type (
	Target struct {
		Host      string
		Path      string
		Port      int32
		IpAddress net.IP
	}

	Route struct {
		Target      *Target
		IngressPath IngressPath
		ServicePort ServicePort
	}
)

func (r *Route) ParseTarget(path string) (err error) {
	r.Target = &Target{}
	var parsedUrl *url.URL
	parsedUrl, err = urlx.Parse(path)
	if err != nil {
		return
	}
	var normalized string
	normalized, err = urlx.Normalize(parsedUrl)
	if err != nil {
		return
	}
	var normalizedUrl *url.URL
	normalizedUrl, err = urlx.Parse(normalized)
	if err != nil {
		return
	}
	var host string
	var port string
	host, port, err = urlx.SplitHostPort(normalizedUrl)
	if err != nil {
		return
	}
	r.Target.Host = host
	if normalizedUrl.Path == "" {
		r.Target.Path = "/"
	} else {
		r.Target.Path = normalizedUrl.Path
	}
	if port != "" {
		var parsedInt int64
		parsedInt, err = strconv.ParseInt(port, 0, 32)
		if err != nil {
			return
		}
		r.Target.Port = int32(parsedInt)
	} else {
		switch normalizedUrl.Scheme {
		case "https":
			r.Target.Port = 443
		case "http":
			r.Target.Port = 80
		}
	}
	var ip *net.IPAddr
	ip, err = urlx.Resolve(normalizedUrl)
	if err != nil {
		return
	}
	r.Target.IpAddress = ip.IP
	return
}

func (r *Route) FindRoute(c *Components) {
	var err error
	var ingressPath *IngressPath
	var podPort *PodPort
	var servicePort *ServicePort
	ingressPath, _ = c.FindIngressPathForHost(r.Target)
	if ingressPath == nil {
		servicePort, err = c.FindServicePortForHost(r.Target)
		if err != nil {
			_ = level.Error(Logger).Log("msg", err)
			os.Exit(1)
		}
		if servicePort == nil {
			_ = level.Error(Logger).Log("msg", "Could not find ingress or service matching host")
			os.Exit(1)
		}
	} else {
		servicePort, err = c.FindServicePortForIngressPath(ingressPath)
		if err != nil {
			_ = level.Error(Logger).Log("msg", err)
			os.Exit(1)
		}
		if servicePort == nil {
			_ = level.Error(Logger).Log("msg", "Could not find service matching ingress rule")
			os.Exit(1)
		}
	}
	podPort, err = c.FindPodPortForServicePort(servicePort)
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
		os.Exit(1)
	}
	if podPort == nil {
		_ = level.Error(Logger).Log("msg", "Could not find pod with port matching service")
		os.Exit(1)
	}
	if ingressPath != nil {
		fmt.Printf(
			"%s:%d%s\n -> ingress: %s\n    path: %s\n    -> service: %s\n       mapping: %s %d -> %s %d\n       -> pod: %s\n          container: %s\n          port: %d\n",
			r.Target.Host,
			r.Target.Port,
			r.Target.Path,
			ingressPath.IngressName,
			ingressPath.Path,
			servicePort.ServiceName,
			servicePort.PortName,
			servicePort.Port,
			servicePort.TargetStrPort,
			servicePort.TargetIntPort,
			podPort.PodName,
			podPort.ContainerName,
			podPort.ContainerPort,
		)
	} else {
		fmt.Printf(
			"%s:%d%s\n -> service: %s\n    mapping: %s %d -> %s %d\n    -> pod: %s\n       container: %s\n       port: %d\n",
			r.Target.Host,
			r.Target.Port,
			r.Target.Path,
			servicePort.ServiceName,
			servicePort.PortName,
			servicePort.Port,
			servicePort.TargetStrPort,
			servicePort.TargetIntPort,
			podPort.PodName,
			podPort.ContainerName,
			podPort.ContainerPort,
		)
	}

}
