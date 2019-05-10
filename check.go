package netkat

import (
	"github.com/go-kit/kit/log/level"
	"github.com/goware/urlx"
	"net"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
)

type (
	Checker struct {
		Target               *Target
		KubernetesRoute      *KubernetesRoute
		KubernetesComponents *KubernetesComponents
		RequiredChecks       []string
		PassedChecks         []string
		FailedChecks         []string
	}

	Check struct {
		Name     string
		Priority int
	}

	KubernetesRoute struct {
		IngressPath
		ServicePort
	}

	Target struct {
		Host      string
		Path      string
		Port      int32
		IpAddress net.IP
	}
)

var (
	checkList = []Check{
		{"CheckKubernetesRouteFromHost", 0},
	}
)

func (ch *Checker) ParseTarget(path string) (err error) {
	ch.Target = &Target{}
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
	ch.Target.Host = host
	if normalizedUrl.Path == "" {
		ch.Target.Path = "/"
	} else {
		ch.Target.Path = normalizedUrl.Path
	}
	if port != "" {
		var parsedInt int64
		parsedInt, err = strconv.ParseInt(port, 0, 32)
		if err != nil {
			return
		}
		ch.Target.Port = int32(parsedInt)
	} else {
		switch normalizedUrl.Scheme {
		case "https":
			ch.Target.Port = 443
		case "http":
			ch.Target.Port = 80
		}
	}
	var ip *net.IPAddr
	ip, err = urlx.Resolve(normalizedUrl)
	if err != nil {
		return
	}
	ch.Target.IpAddress = ip.IP
	return
}

func (ch *Checker) RunChecks() {
	ch.InitChecks()
	for _, check := range ch.RequiredChecks {
		reflect.ValueOf(ch).MethodByName(check).Call([]reflect.Value{})
	}
	PrintCheckResults(ch)
}

func (ch *Checker) InitChecks() {
	sort.Slice(checkList, func(i, j int) bool { return checkList[i].Priority < checkList[j].Priority })
	for _, check := range checkList {
		ch.RequiredChecks = append(ch.RequiredChecks, check.Name)
	}
}

func (ch *Checker) CheckKubernetesRouteFromHost() {
	PrintCheckHeader()
	var err error
	var ingressPath *IngressPath
	var podPorts []*PodPort
	var servicePort *ServicePort
	indent := 0
	PrintHost(ch.Target)
	ingressPath, _ = ch.KubernetesComponents.FindIngressPathForHost(ch.Target)
	if ingressPath == nil {
		servicePort, err = ch.KubernetesComponents.FindServicePortForHost(ch.Target)
		if err != nil {
			_ = level.Error(Logger).Log("msg", err)
			os.Exit(1)
		}
		if servicePort == nil {
			_ = level.Error(Logger).Log("msg", "Could not find ingress or service matching host")
			os.Exit(1)
		}
		PrintServicePort(servicePort, indent)
		indent = indent + 3
	} else {
		PrintIngressPath(ingressPath, indent)
		indent = indent + 3
		servicePort, err = ch.KubernetesComponents.FindServicePortForIngressPath(ingressPath)
		if err != nil {
			_ = level.Error(Logger).Log("msg", err)
			os.Exit(1)
		}
		if servicePort == nil {
			_ = level.Error(Logger).Log("msg", "Could not find service matching ingress rule")
			os.Exit(1)
		}
		PrintServicePort(servicePort, indent)
		indent = indent + 3
	}
	podPorts, err = ch.KubernetesComponents.FindPodPortForServicePort(servicePort)
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
		os.Exit(1)
	}
	for _, p := range podPorts {
		PrintPodPort(p, indent)
	}
	PrintPassFooter(ch)
}
