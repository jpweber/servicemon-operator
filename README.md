# Service Monitor Operator

## Overview

The Service Monitor Operator will create, destroy and update service monitor resources for you when you create, update or destroy a service.

## Why

The prometheus operator leverages `Service Monitors` Custom Resrouces to define what services to monitor and details around how to scrape them. The goal of this project is to _try_ and solve the most common service monitor configurations by automatically creating them for you based on annotations on your service.

## Use

The Service Monitor Operator will look at the following annotations on a service.

* `prometheus.io/probe:` - true or false
* `prometheus.io/port:` - TCP port listening for scrape requests
* `prometheus.io/path:` - path to metrics endpoint
* `prometheus.io/scheme:` - http or https
* `prometheus.io/interval:`- scrape interval
* `prometheus.io/bearertoken:` path to bearer token file

These are all optional except for `prometheus.io/probe`. If any values are omitted, the following are the default values that will be used

* port - `8080`
* path - `/metrics`
* scheme - `http`
* scrapeInterval - `30s`
* bearerTokenFile - `/var/run/secrets/kubernetes.io/serviceaccount/token`

If the `prometheus.io/probe:` annotation valuel is _not_ set to `true` it will ignore this service.  
The path annotation is the path to the endpoint that will provide metrics. For example in the case of `foo.bar.com/metrics` the path is `/metrics`

The Selectors for the `servicemonitor` will be set to match the `Labels` on the service and the `Namespace` that the service is running in.

### Unimplmented parts of the ServiceMonitor spec

The following items are not _yet_ implmented

* TargetLabels are not able to be changed
* PodTargetLabels are not able to be changed
* TLSConfig is not supported

## Deploy to cluster

`kubectl apply -f deploy -n monitoring`

## Build

* Requires operator-sdk from the [operator framework](https://github.com/operator-framework/operator-sdk)
* To build an image for deployment/testing run the following
  * `operator-sdk build <docker image name of your choosing>`
* Update `deploy/operator.yaml` file with your image name
* Deploy to cluster `kubectl apply -f deploy -n monitoring`

You can also test without building and deploying to a cluster with the following command
`operator-sdk up local --namespace=<your namespace> --kubeconfig=<path to your kubeconfig file>`

You can ommit the `kubeconfig` flag if you are using the default kubeconfig on your system `~/.kube/config`.