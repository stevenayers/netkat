package netkat_test

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/stevenayers/netkat"
	"github.com/stretchr/testify/suite"
	"os"
	"os/user"
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
	usr, _ := user.Current()
	netkat.InitLogger(log.NewSyncWriter(os.Stdout), "error")
	s.client = netkat.InitClient("minikube", fmt.Sprintf("%v/.kube/config", usr.HomeDir))
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
