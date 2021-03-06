# netkat
[![Build Status](https://travis-ci.org/stevenayers/netkat.svg?branch=master)](https://travis-ci.org/stevenayers/netkat)
[![codecov.io Code Coverage](https://img.shields.io/codecov/c/github/stevenayers/netkat.svg)](https://codecov.io/github/stevenayers/netkat?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/stevenayers/netkat)](https://goreportcard.com/report/github.com/stevenayers/netkat)
[![Release](https://img.shields.io/badge/release-v0.1--alpha-5272B4.svg)](https://github.com/stevenayers/netkat/releases/tag/v0.1-alpha)
[![GoDoc](https://godoc.org/github.com/stevenayers/netkat?status.svg)](https://godoc.org/github.com/stevenayers/netkat)

CLI for troubleshooting kubernetes networking issues.

## Getting Started

Build from source:
* Requires Go 1.13 and dep package management
```bash
git clone git@github.com:stevenayers/netkat.git
cd netkat
go build
go build cmd/main.go
mv ./main /usr/local/bin/netkat
```
For help:
```bash
$ netkat -h
```

Example Usage:
```bash
$ netkat grafana.digital.foobar.com -context kops-dev -config ~/.kube/config
$ netkat pod/grafana-fb86ad62c-f63x9:3000 -context kops-dev -config ~/.kube/config
```
```
=== RUN   CheckKubernetesRouteFromHost
host: grafana.digital.foobar.com
port: 80
path: /
ip address: 34.89.100.1
 -> ingress: grafana-ingress
    namespace: metrics
    path: /
    ip address: 34.89.100.1
    -> service: grafana-service
       namespace: metrics
       app selector: grafana-app
       external IP: 34.89.100.1
       internal IP: 10.44.0.1
       mapping: http (80) -> 3000
       -> pod: grafana-fb86ad62c-p72v8
          namespace: metrics
          app: grafana-app
          container: grafana
          port: 3000
       -> pod: grafana-fb86ad62c-lg92a
          namespace: metrics
          app: grafana-app
          container: grafana
          port: 3000
       -> pod: grafana-fb86ad62c-f63x9
          namespace: metrics
          app: grafana-app
          container: grafana
          port: 3000
--- PASS: CheckKubernetesRouteFromHost
=== PASS: (1/1)
    --- CheckKubernetesRouteFromHost
=== FAIL: (0/1)
```

Under development, current version will only print out the route when config is setup correctly.
Incorrect configuration just throws an error and prints out nothing. This needs to be implemented properly.

* Checks ownership of DNS records (to be implemented)
* Checks external DNS logs (to be implemented)
* Matches A record against ingress/service
* Checks service/ingress config
* Checks ports mappings
* Checks port is open on pod
* Checks LB rules on cloud provider side (to be implemented)
* Checks LoadBalancerSourceRanges (to be implemented)


## What Done Looks Like
End-to-end Scenarios
```
local -> pod_name:port
local -> fqdn:port
local -> http(s)://url/path
local -> http(s)://url:port/path
```

|**Check Name**|**Description**|**Done**|
|:-----:|:-----:|:-----:|
CheckKubernetesRouteFromHost| Takes the host:port info and matches it to ingress or/then service then pod. | x
CheckStatusPod|  Checks pod status is running| x
CheckListeningPod|  Portforwards directly to pod and checks connection| x
CheckKubernetesRouteFromPod| Takes pod:port and maps backwards to a hostname then checks the host configuration. | x
CheckKubernetesRouteFromInternalHost| Takes the host:port info and matches it to ingress or/then service then pod but for intra-cluster situations. | 
CheckKubernetesRoutePodToPod| Takes pod:port and maps to pod:port| 
CheckStatusNginxIngress| Checks nginx-ingress is healthy.| 
CheckStatusTraefikIngress| Checks traefik ingress is healthy.| 
CheckStatusKubeDns| Checks kube-dns is healthy.| 
CheckSourceRangesIngress| Checks any source range annotations on ingress against originating IP. | 
CheckSourceRangesService| Checks any source range annotations on service against originating IP. | 
CheckInboundRulesLB| Checks originating IP against inbound rules for Load Balancer. | 
CheckInboundRulesLBAzure|  hecks originating IP against inbound rules for Load Balancer. | 
CheckInboundRulesLBAWS| Checks originating IP against inbound rules for Load Balancer. | 
CheckDnsOwnershipGCP| | 
CheckDnsOWnershipAzure| | 
CheckDnsOwnershipAWS| | 
CheckDnsInternalPodToPod| | 
