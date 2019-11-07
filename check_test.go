package netkat_test

import (
	"fmt"
	"github.com/stevenayers/netkat"
	"github.com/stretchr/testify/assert"
)

type (
	TargetTest struct {
		RawUrl string
		Host   string
		Port   int32
		Path   string
	}

	PodTest struct {
		PodPort  netkat.PodPort
		Expected int
	}

	ServiceTest struct {
		ServicePort netkat.ServicePort
		Expected    int
	}
)

var (
	TargetTests = []TargetTest{
		{"google.com", "google.com", 80, "/"},
		{"https://google.com", "google.com", 443, "/"},
		{"http://google.com", "google.com", 80, "/"},
		{"http://google.com:8000", "google.com", 8000, "/"},
		{"google.com:8000", "google.com", 8000, "/"},
		{"google.com/path", "google.com", 80, "/path"},
		{"https://google.com/path", "google.com", 443, "/path"},
		{"http://google.com/path", "google.com", 80, "/path"},
		{"http://google.com:8000/path", "google.com", 8000, "/path"},
		{"google.com:8000/path", "google.com", 8000, "/path"},
	}
)

func (s *StoreSuite) TestTarget() {
	for _, test := range TargetTests {
		var r netkat.Checker
		err := r.ParseTarget(test.RawUrl)
		if err != nil {
			s.T().Fatal(err)
		}
		assert.Equal(s.T(), test.Host, r.Target.Host)
		assert.Equal(s.T(), test.Port, r.Target.Port)
		assert.Equal(s.T(), test.Path, r.Target.Path)
		fmt.Println(r.Target.IpAddress.String())
	}

}

func (s *StoreSuite) TestRunChecks() {
	var ch netkat.Checker
	err := ch.ParseTarget(s.target)
	if err != nil {
		s.T().Fatal(err)
	}
	ch.KubernetesComponents = s.client.GetComponents()
	ch.Client = s.client
	ch.RunChecks()
	assert.Equal(s.T(), 3, len(ch.PassedChecks), "Expected checks to pass")
}

func (s *StoreSuite) TestCheckKubernetesRouteFromHost() {
	var ch netkat.Checker
	err := ch.ParseTarget(s.target)
	if err != nil {
		s.T().Fatal(err)
	}
	ch.KubernetesComponents = s.client.GetComponents()
	ch.CheckKubernetesRouteFromHost()
	assert.Equal(s.T(), 1, len(ch.PassedChecks), "Expected CheckKubernetesRouteFromHost to pass")
}

func (s *StoreSuite) TestCheckStatusPod() {
	var ch netkat.Checker
	err := ch.ParseTarget(s.target)
	if err != nil {
		s.T().Fatal(err)
	}
	ch.KubernetesComponents = s.client.GetComponents()
	ch.KubernetesRoute = &netkat.KubernetesRoute{}
	ch.KubernetesRoute.Pods = ch.KubernetesComponents.PodPorts
	ch.CheckStatusPod()
	assert.Equal(s.T(), 1, len(ch.PassedChecks), "Expected CheckStatusPod to pass")
}

func (s *StoreSuite) TestCheckListeningPod() {
	var ch netkat.Checker
	ch.KubernetesComponents = s.client.GetComponents()
	PodTests := []PodTest{
		{netkat.PodPort{PodName: ch.KubernetesComponents.PodPorts[0].PodName, Namespace: "default", ContainerPort: 8080}, 1},
		{netkat.PodPort{PodName: ch.KubernetesComponents.PodPorts[0].PodName, Namespace: "default", ContainerPort: 54921}, 0},
		{netkat.PodPort{PodName: "bad-name", Namespace: "default", ContainerPort: 8080}, 0},
	}
	for _, test := range PodTests {
		var ch netkat.Checker
		err := ch.ParseTarget(s.target)
		if err != nil {
			s.T().Fatal(err)
		}
		ch.Client = s.client
		ch.KubernetesRoute = &netkat.KubernetesRoute{}
		ch.KubernetesRoute.Pods = []*netkat.PodPort{&test.PodPort}
		ch.CheckListeningPod()
		assert.Equal(s.T(), test.Expected, len(ch.PassedChecks), "Expected CheckListeningPod to pass",
			test.PodPort.PodName, test.PodPort.ContainerPort)
	}
}
