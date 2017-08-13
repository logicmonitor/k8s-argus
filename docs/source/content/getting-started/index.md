---
title: "Getting Started"
date: 2017-08-12T16:20:39-07:00
draft: false
---

{{< warning title="Warning" >}}
Argus is a community driven project. LogicMonitor support will not assist in any issues related to Argus.
{{< /warning >}}

[Argus](https://github.com/logicmonitor/k8s-argus) can be installed with [Helm](https://github.com/kubernetes/helm). A [LogicMonitor](https://www.logicmonitor.com) account is required.

First, you will need to add the LogicMonitor chart repository using the Helm CLI:
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

Required Values:

-   **accessID:** The LogicMonitor API key ID.
-   **accessKey:** The LogicMonitor API key.
-   **account:** The LogicMonitor account name.
-   **clusterName:** A unique name given to the cluster's device group.
-   **collectorVersion:** The collector version to install.
-   **collectorDescription:** A unique collector description used to look up a collector dynamically.

Optional Values:

-   **collectorImageTag:** The collector image tag.
-   **collectorSize:** The collector size to install. Can be nano, small, medium, or large.
-   **etcdDiscoveryToken:** The public etcd discovery token used to add etcd hosts to the cluster device group.
-   **imageTag:** The argus image tag to use.
