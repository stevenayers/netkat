# netkat
CLI for troubleshooting kubernetes networking issues.

Example Usage:
```bash
$ netkat -target grafana.digital.foobar.com -context kops-dev -config ~/.kube/config
grafana.digital.foobar.com:80/
 -> ingress: grafana-ingress
    path: /
    -> service: grafana-service
       mapping: http 80 -> 3000
       -> pod: grafana-fb86ad62c-p72v8
          container: grafana
          port: 3000
       -> pod: grafana-fb86ad62c-lg92a
          container: grafana
          port: 3000
       -> pod: grafana-fb86ad62c-f63x9
          container: grafana
          port: 3000
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



