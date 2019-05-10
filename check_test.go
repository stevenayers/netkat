package netkat_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"netkat"
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

func (s *StoreSuite) TestCheckKubernetesRouteFromHost() {
	var ch netkat.Checker
	var targetString string
	err := ch.ParseTarget(targetString)
	if err != nil {
		s.T().Fatal(err)
	}
	ch.KubernetesComponents = s.components
	ch.RunChecks()
}
