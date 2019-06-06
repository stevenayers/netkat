package netkat_test

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *StoreSuite) TestGetPods() {
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
	pods, err := s.client.CoreV1().Services("").List(v1.ListOptions{})
	if err != nil {
		s.T().Fatal(err)
	}
	for _, pod := range pods.Items {
		fmt.Print(pod)
	}
}

func (s *StoreSuite) TestGetIngress() {
	pods, err := s.client.ExtensionsV1beta1().Ingresses("").List(v1.ListOptions{})
	if err != nil {
		s.T().Fatal(err)
	}
	for _, pod := range pods.Items {
		fmt.Print(pod)
	}
}

func (s *StoreSuite) TestGetDeployments() {
	pods, err := s.client.ExtensionsV1beta1().Deployments("").List(v1.ListOptions{})
	if err != nil {
		s.T().Fatal(err)
	}
	for _, pod := range pods.Items {
		fmt.Print(pod)
	}
}
