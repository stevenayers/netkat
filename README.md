# netkat
CLI for troubleshooting kubernetes networking issues.

## Getting Started

Build from source:
* Requires Go 1.11.5 and dep package management
```bash
git clone git@github.com:stevenayers/netkat.git
cd netkat
dep ensure
go build cmd/main.go
mv ./main /usr/local/bin/netkat
```
For help:
```bash
$ netkat -h
```

Example Usage:
```bash
$ netkat -target grafana.digital.foobar.com -context kops-dev -config ~/.kube/config
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
* Checks port is open on pod (to be implemented)
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

Checks
```
- CheckKubernetesRouteFromHost
  Takes the host:port info and matches it to ingress, or/then service, then pod.
- CheckKubernetesRouteFromPod (Not yet implemented)
  Takes pod:port and maps backwards to a hostname, then checks the host configuration.
- CheckKubernetesRouteFromInternalHost (Not yet implemented)
  Takes the host:port info and matches it to ingress, or/then service, then pod, but for intra-cluster situations.
- CheckKubernetesRoutePodToPod (Not yet implemented)
  Takes pod:port and maps to pod:port
- CheckStatusPod (Not yet implemented)
  Checks pod status is running
- CheckStatusIngressController (Not yet implemented)
  Checks nginx-ingress is healthy.
- CheckStatusKubeDns (Not yet implemented)
  Checks kube-dns is health
- CheckListeningPod (Not yet implemented)
  Portforwards directly to pod and checks connection
- CheckListeningService (Not yet implemented)
  Portforwards to service and checks connection
- CheckSourceRangesIngress (Not yet implemented)
  Checks any source range annotations on ingress against originating IP
- CheckSourceRangesService (Not yet implemented)
  Checks any source range annotations on service against originating IP
- CheckInboundRulesLBGCP (Not yet implemented)
  Checks originating IP against inbound rules for Load Balancer
- CheckInboundRulesLBAzure (Not yet implemented)
  Checks originating IP against inbound rules for Load Balancer
- CheckInboundRulesLBAWS (Not yet implemented)
  Checks originating IP against inbound rules for Load Balancer
- CheckDnsOwnershipGCP (Not yet implemented)
- CheckDnsOWnershipAzure (Not yet implemented)
- CheckDnsOwnershipAWS (Not yet implemented)
- CheckDnsInternalPodToPod (Not yet implemented)
```

