---
title: "Getting Started"
date: 2017-08-12T16:20:39-07:00
draft: false
---

# Quick Start
The simplest way to install Argus is to use [Helm](https://github.com/kubernetes/helm). First, you will need to add the LogicMonitor chart repository:
```bash
$ helm repo add logicmonitor https://logicmonitor.github.com/k8s-helm-charts
```

Now, install Argus:
```bash
$ helm upgrade
    --install \
    --debug \
    --wait \
    --namespace '$NAMESPACE' \
    --set accessID='$ACCESS_ID' \
    --set accessKey='$ACCESS_KEY' \
    --set account='$ACCOUNT' \
    --set clusterName='$CLUSTER_NAME' \
    --set collectorDescription='$COLLECTOR_DESCRIPTION' \
    --set collectorImageTag='$COLLECTOR_IMAGE_TAG' \
    --set collectorSize='$COLLECTOR_SIZE' \
    --set collectorVersion='$COLLECTOR_VERSION' \
    --set etcdDiscoveryToken='$ETCD_DISCOVERY_TOKEN' \
    --set imageTag='$IMAGE_TAG' \
    argus logicmonitor/argus
```
> Note: Argus should be installed only once per cluster.

# Community

-   To report bugs and/or submit feature requests, use [GitHub](https://github.com/logicmonitor/k8s-argus/issues).
