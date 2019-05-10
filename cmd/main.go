package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	"netkat"
	"os"
)

var (
	target  = flag.String("target", "github.com:443/uri", "Specify target host and port")
	context = flag.String("context", "default", "Specify kubernetes cluster context of pod")
	config  = flag.String("config", "./config", "Kubernetes config file")
)

func main() {
	flag.Parse()
	var ch netkat.Checker
	netkat.InitLogger(log.NewSyncWriter(os.Stdout), "error")
	client := netkat.InitClient(*context, *config)

	err := ch.ParseTarget(*target)
	if err != nil {

	}
	ch.KubernetesComponents = client.GetComponents()
	ch.CheckKubernetesRouteFromHost()

}
