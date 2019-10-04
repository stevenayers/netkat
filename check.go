package netkat

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/goware/urlx"
	"net"
	"net/url"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type (
	Checker struct {
		Target               *Target
		KubernetesRoute      *KubernetesRoute
		KubernetesComponents *KubernetesComponents
		Client               Client
		RequiredChecks       []string
		PassedChecks         []string
		FailedChecks         []string
	}

	Check struct {
		Name     string
		Priority int
	}

	KubernetesRoute struct {
		Ingress *IngressPath
		Service *ServicePort
		Pods    []*PodPort
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
		{"CheckStatusPod", 1},
		{"CheckListeningPod", 2},
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

func (ch *Checker) PassCheck() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fullName := strings.Split(frame.Function, ".")
	functionName := fullName[len(fullName)-1]
	fmt.Printf(
		"--- PASS: %s\n", functionName,
	)
	ch.PassedChecks = append(ch.PassedChecks, functionName)
}

func (ch *Checker) FailCheck() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fullName := strings.Split(frame.Function, ".")
	functionName := fullName[len(fullName)-1]
	fmt.Printf(
		"--- FAIL: %s\n", functionName,
	)
	ch.PassedChecks = append(ch.FailedChecks, functionName)
}

func (ch *Checker) CheckKubernetesRouteFromHost() {
	PrintCheckHeader()
	var err error
	indent := 0
	ch.KubernetesRoute = &KubernetesRoute{}
	PrintHost(ch.Target)
	ch.KubernetesRoute.Ingress, _ = ch.KubernetesComponents.FindIngressPathForHost(ch.Target)
	if ch.KubernetesRoute.Ingress == nil {
		ch.KubernetesRoute.Service, err = ch.KubernetesComponents.FindServicePortForHost(ch.Target)
		if err != nil {
			_ = level.Error(Logger).Log("msg", err)
			ch.FailCheck()
			return
		}
		if ch.KubernetesRoute.Service == nil {
			_ = level.Error(Logger).Log("msg", "Could not find ingress or service matching host")
			ch.FailCheck()
			return
		}
		PrintServicePort(ch.KubernetesRoute.Service, indent)
		indent = indent + 3
	} else {
		PrintIngressPath(ch.KubernetesRoute.Ingress, indent)
		indent = indent + 3
		ch.KubernetesRoute.Service, err = ch.KubernetesComponents.FindServicePortForIngressPath(ch.KubernetesRoute.Ingress)
		if err != nil {
			_ = level.Error(Logger).Log("msg", err)
			ch.FailCheck()
			return
		}
		if ch.KubernetesRoute.Service == nil {
			_ = level.Error(Logger).Log("msg", "Could not find service matching ingress rule")
			ch.FailCheck()
			return
		}
		PrintServicePort(ch.KubernetesRoute.Service, indent)
		indent = indent + 3
	}
	ch.KubernetesRoute.Pods, err = ch.KubernetesComponents.FindPodPortForServicePort(ch.KubernetesRoute.Service)
	if err != nil {
		_ = level.Error(Logger).Log("msg", err)
		ch.FailCheck()
		return
	}
	for _, p := range ch.KubernetesRoute.Pods {
		PrintPodPort(p, indent)
	}
	ch.PassCheck()
}

func (ch *Checker) CheckStatusPod() {
	PrintCheckHeader()
	if len(ch.KubernetesRoute.Pods) > 0 {
		for _, p := range ch.KubernetesRoute.Pods {
			if p.PodStatus != "Running" {
				_ = level.Error(Logger).Log("msg", "Not all pods have a status of `Running`.")
				ch.FailCheck()
				return
			}
		}
	} else {
		_ = level.Error(Logger).Log("msg", "No pods were found.")
		ch.FailCheck()
		return
	}
	ch.PassCheck()
}

func (ch *Checker) CheckListeningPod() {
	PrintCheckHeader()
	if len(ch.KubernetesRoute.Pods) > 0 {
		for _, p := range ch.KubernetesRoute.Pods {
			res, err := ch.Client.GetPortforwardResponse(p)
			if err != nil {
				_ = level.Error(Logger).Log(
					"msg", fmt.Sprintf("Error connecting to port: %v", err))
				return
			}
			if res.StatusCode != 200 {
				_ = level.Error(Logger).Log(
					"msg", fmt.Sprintf("Bad HTTP Status Code: %v", res.StatusCode))
				ch.FailCheck()
				return
			}
		}
	} else {
		_ = level.Error(Logger).Log("msg", "No pods were found.")
		ch.FailCheck()
		return
	}
	ch.PassCheck()
}
