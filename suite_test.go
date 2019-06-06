package netkat_test

import (
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/suite"
	"netkat"
	"os"
	"testing"
)

type (
	StoreSuite struct {
		suite.Suite
		client     netkat.Client
		components *netkat.KubernetesComponents
		target     string
	}
)

func (s *StoreSuite) SetupSuite() {
	netkat.InitLogger(log.NewSyncWriter(os.Stdout), "error")
	s.client = netkat.InitClient("netkat-test", "/Users/stevenayers/.kube/config")
	s.target = "http://hello-world.info"
}

func (s *StoreSuite) SetupTest() {

}

func (s *StoreSuite) TearDownSuite() {

}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}
