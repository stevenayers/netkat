package netkat_test

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/suite"
	"netkat"
	"os"
	"testing"
)

type (
	StoreSuite struct {
		suite.Suite
		client netkat.Client
		driver neo4j.Driver
	}
)

func (s *StoreSuite) SetupSuite() {
	netkat.InitLogger(log.NewSyncWriter(os.Stdout), "error")
	s.client = netkat.InitClient("default", "./config")
}

func (s *StoreSuite) SetupTest() {

}

func (s *StoreSuite) TearDownSuite() {
	err := s.driver.Close()
	if err != nil {
		fmt.Print(err)
	}
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}
