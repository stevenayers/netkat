package netkat_test

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"netkat"
	"os"
	"testing"
)

type (
	StoreSuite struct {
		suite.Suite
		client     netkat.Client
		components *netkat.KubernetesComponents
	}
)

func (s *StoreSuite) SetupSuite() {
	netkat.InitLogger(log.NewSyncWriter(os.Stdout), "error")
	s.client = netkat.InitClient("default", "./config")
}

func (s *StoreSuite) SetupTest() {

}

func (s *StoreSuite) TearDownSuite() {

}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}
