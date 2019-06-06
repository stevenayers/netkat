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
	ch.RunChecks()
	assert.Equal(s.T(), 2, len(ch.PassedChecks), "Expected checks to pass")
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
