# Service Monitor Operator

## Overview
The Service Monitor Operator will create, destroy and update service monitor resources for you when you create, update or destroy a service. 

## Why
The prometheus operator leverages `Service Monitors` Custom Resrouces to define what services to monitor and details around how to scrape them. The goal of this project is to _try_ and solve the most common service monitor configurations by automatically creating them for you based on annotations on your service.

## Use
The Service Monitor Operator looks for the following annotations on a service.
* `prometheus.io/probe: ` - true or false
* `prometheus.io/path: ` - path to metrics endpoint
* `prometheus.io/port: ` - TCP port listening for scrape requests

If the `prometheus.io/probe:` annotation valuel is not set to `true` it will ignore this service.  
The path annotation is the path to the endpoint that will provide metrics. For example in the case of `foo.bar.com/metrics` the path is `/metrics`

The Selectors for the `servicemonitor` will be set to match the `Labels` on the service and the `Namespace` that the service is running in. 

### Unimplmented parts of the ServiceMonitor spec
The following items are not _yet_ implmented
* TargetLabels are not able to be changed
* PodTargetLabels are not able to be changed
* Scheme is hardcoded to http
* TLSConfig is not supported
* Changing scrape interval
* Changing `SampleLimit`

## Install
`kubectl apply -f deploy -n monitoring`