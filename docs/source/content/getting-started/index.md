---
title: "Getting Started"
date: 2017-08-12T16:20:39-07:00
draft: false
---

# Quick Start

The simplest way to install Argus is to use [Helm](https://github.com/kubernetes/helm). First, you will need to add the LogicMonitor chart repository:

```bash
$ helm repo add logicmonitor https://logicmonitor.github.com/k8s-helm-charts
"logicmonitor" has been added to your repositories
```

Next, install the LogicMonitor Collectorset controller:

```bash
$ helm upgrade \
  --install \
  --debug \
  --wait \
  --tiller-namespace="$NAMESPACE" \
  --set accessID="$ACCESS_ID" \
  --set accessKey="$ACCESS_KEY" \
  --set account="$ACCOUNT" \
  --set clusterName="$CLUSTER_NAME" \
  --set etcdDiscoveryToken="$ETCD_DISCOVERY_TOKEN" \
  --set imageTag="$IMAGE_TAG" \
  collectorset-controller logicmonitor/collectorset-controller
```

> Note: The Collectorset controller should be installed only once per cluster.

Now, install Argus:

```bash
    $ helm upgrade \
    --install \
    --debug \
    --wait \
    --tiller-namespace="$NAMESPACE" \
    --set accessID="$ACCESS_ID" \
    --set accessKey="$ACCESS_KEY" \
    --set account="$ACCOUNT" \
    --set clusterName="$CLUSTER_NAME" \
    --set etcdDiscoveryToken="$ETCD_DISCOVERY_TOKEN" \
    --set imageTag="$IMAGE_TAG" \
    --set collector.replicas="$COLLECTOR_REPLICAS" \
    --set collector.size="$COLLECTOR_SIZE" \
    argus logicmonitor/argus
```

> Note: Argus should be installed only once per cluster.

# Community

- To report bugs and/or submit feature requests, use [GitHub](https://github.com/logicmonitor/k8s-argus/issues).
