package netkat_test

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *StoreSuite) TestGets() {
	pods := s.client.GetPods()
	services := s.client.GetServices()
	ingresses := s.client.GetIngresses()
	po, err := json.Marshal(pods)
	if err != nil {
		fmt.Print(po)
	}
	se, err := json.Marshal(services)
	if err != nil {
		fmt.Print(se)
	}
	in, err := json.Marshal(ingresses)
	if err != nil {
		fmt.Print(in)
	}

}

func (s *StoreSuite) TestGetServices() {
	_, err := s.client.CoreV1().Services("").List(v1.ListOptions{})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TestGetIngress() {
	_, err := s.client.ExtensionsV1beta1().Ingresses("").List(v1.ListOptions{})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TestGetPods() {
	_, err := s.client.CoreV1().Pods("").List(v1.ListOptions{})
	if err != nil {
		s.T().Fatal(err)
	}
}
