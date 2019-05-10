# netkat
CLI for troubleshooting kubernetes networking issues.

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
local -> hostname:port
local -> http(s)://url/path
local -> http(s)://url:port/path
```

Checks
```
 - CheckKubernetesRouteFromHost
 - CheckKubernetesRouteFromPod (Not yet implemented)
 - CheckKubernetesInternalRouteToPod (Not yet implemented)
 - CheckKubernetesRoutePodToPod (Not yet implemented)
 - CheckStatusPod (Not yet implemented)
 - CheckStatusIngressController (Not yet implemented)
 - CheckStatusKubeDns (Not yet implemented)
 - CheckListeningPod (Not yet implemented)
 - CheckListeningService (Not yet implemented)
 - CheckSourceRangesIngress (Not yet implemented)
 - CheckSourceRangesService (Not yet implemented)
 - CheckInboundRulesLBGCP (Not yet implemented)
 - CheckInboundRulesLBAzure (Not yet implemented)
 - CheckInboundRulesLBAWS (Not yet implemented)
 - CheckDnsOwnershipGCP (Not yet implemented)
 - CheckDnsOWnershipAzure (Not yet implemented)
 - CheckDnsOwnershipAWS (Not yet implemented)
 - CheckDnsInternalPodToPod (Not yet implemented)
```

