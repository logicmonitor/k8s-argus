---
title: "Getting Started"
date: 2017-08-12T16:20:39-07:00
draft: false
---

The simplest way to get started with Argus is to install it using [Helm]
(https://github.com/kubernetes/helm). Prior to installation, you will need a
cluster-admin serviceaccount for tiller:
```bash
$ kubectl create serviceaccount tiller --namespace $NAMESPACE
$ kubectl create clusterrolebinding tiller --clusterrole=cluster-admin --serviceaccount=$NAMESPACE:tiller
$ helm init --service-account=tiller
```

You'll also need to add the LogicMonitor chart repository:

```bash
$ helm repo add logicmonitor https://logicmonitor.github.com/k8s-helm-charts
"logicmonitor" has been added to your repositories
```

Now you can install the LogicMonitor Collectorset controller:

```bash
$ helm upgrade \
  --install \
  --debug \
  --wait \
  --set accessID="$ACCESS_ID" \
  --set accessKey="$ACCESS_KEY" \
  --set account="$ACCOUNT" \
  --set clusterName="$CLUSTER_NAME" \
  --set etcdDiscoveryToken="$ETCD_DISCOVERY_TOKEN" \
  --set imageTag="$IMAGE_TAG" \
  collectorset-controller logicmonitor/collectorset-controller
```

> Note: The Collectorset controller should be installed only once per cluster.

Next, install Argus:

```bash
    $ helm upgrade \
    --install \
    --debug \
    --wait \
    --set accessID="$ACCESS_ID" \
    --set accessKey="$ACCESS_KEY" \
    --set account="$ACCOUNT" \
    --set clusterGroupID="$CLUSTER_GROUP_ID" \
    --set clusterName="$CLUSTER_NAME" \
    --set etcdDiscoveryToken="$ETCD_DISCOVERY_TOKEN" \
    --set imageTag="$IMAGE_TAG" \
    --set collector.replicas="$COLLECTOR_REPLICAS" \
    --set collector.size="$COLLECTOR_SIZE" \
    argus logicmonitor/argus
```

> Note: Argus should be installed only once per cluster.

After installation is complete, you should make sure you have [DataSources]
(https://logicmonitor.github.io/k8s-argus/docs/monitoring/) in your account
that will start monitoring the resources in your cluster.

# Community

- To report bugs and/or submit feature requests, use [GitHub]
(https://github.com/logicmonitor/k8s-argus/issues).
