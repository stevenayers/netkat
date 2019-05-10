package netkat

import (
	"fmt"
	"runtime"
	"strings"
)

func PrintCheckHeader() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fullName := strings.Split(frame.Function, ".")
	functionName := fullName[len(fullName)-1]
	fmt.Printf(
		"=== RUN   %s\n", functionName,
	)
}

func PrintCheckResults(ch *Checker) {

	fmt.Printf("=== PASS: (%d/%d)\n", len(ch.PassedChecks), len(ch.RequiredChecks))
	for _, functionName := range ch.PassedChecks {
		fmt.Printf(
			"    --- %s\n", functionName,
		)
	}
	fmt.Printf("=== FAIL: (%d/%d)\n", len(ch.FailedChecks), len(ch.RequiredChecks))
	for _, functionName := range ch.FailedChecks {
		fmt.Printf(
			"    --- %s\n", functionName,
		)
	}

}

func PrintPassFooter(ch *Checker) {
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

func PrintHost(t *Target) {
	fmt.Printf("host: %s\n", t.Host)
	fmt.Printf("port: %d\n", t.Port)
	fmt.Printf("path: %s\n", t.Path)
	fmt.Printf("ip address: %s\n", t.IpAddress)

}

func PrintIngressPath(i *IngressPath, indent int) {
	fmt.Printf("%v-> ingress: %s\n", strings.Repeat(" ", indent), i.IngressName)
	fmt.Printf("%v   namespace: %s\n", strings.Repeat(" ", indent), i.Namespace)
	fmt.Printf("%v   path: %s\n", strings.Repeat(" ", indent), i.IngressName)
	fmt.Printf("%v   ip address: %s\n", strings.Repeat(" ", indent), i.IpAddress)
}

func PrintServicePort(s *ServicePort, indent int) {
	var srcPort string
	var dstPort string
	if s.SourcePortName != "" {
		srcPort = fmt.Sprintf("%s (%d)", s.SourcePortName, s.SourcePort)
	} else if s.SourcePort == 0 {
		srcPort = fmt.Sprintf("%s", s.SourcePortName)
	} else {
		srcPort = fmt.Sprintf("%d", s.SourcePort)
	}

	if s.TargetPortName != "" {
		dstPort = fmt.Sprintf("%s (%d)", s.TargetPortName, s.TargetPort)
	} else if s.TargetPort == 0 {
		dstPort = fmt.Sprintf("%s", s.TargetPortName)
	} else {
		dstPort = fmt.Sprintf("%d", s.TargetPort)
	}

	fmt.Printf("%v-> service: %s\n", strings.Repeat(" ", indent), s.ServiceName)
	fmt.Printf("%v   namespace: %s\n", strings.Repeat(" ", indent), s.Namespace)
	fmt.Printf("%v   app selector: %s\n", strings.Repeat(" ", indent), s.AppSelector)
	fmt.Printf("%v   external IP: %s\n", strings.Repeat(" ", indent), s.ExternalIP)
	fmt.Printf("%v   internal IP: %s\n", strings.Repeat(" ", indent), s.ClusterIP)
	fmt.Printf("%v   mapping: %s -> %s\n", strings.Repeat(" ", indent), srcPort, dstPort)
}

func PrintPodPort(p *PodPort, indent int) {
	fmt.Printf("%v-> pod: %s\n", strings.Repeat(" ", indent), p.PodName)
	fmt.Printf("%v   namespace: %s\n", strings.Repeat(" ", indent), p.Namespace)
	fmt.Printf("%v   app: %s\n", strings.Repeat(" ", indent), p.App)
	fmt.Printf("%v   container: %s\n", strings.Repeat(" ", indent), p.ContainerName)
	fmt.Printf("%v   port: %d\n", strings.Repeat(" ", indent), p.ContainerPort)
}
