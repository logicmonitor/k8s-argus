---
title: "Getting Started"
date: 2017-08-12T16:20:39-07:00
draft: false
---

The simplest way to get started with Argus is to install it using [Helm](https://helm.sh/docs) version 3.

You'll need to add the LogicMonitor chart repository:

```bash
$ helm repo add logicmonitor https://logicmonitor.github.com/k8s-helm-charts
```

> Note: Argus helm charts will only be installed using Helm 3 on Kubernetes clusters newer than version 1.14.0. For any reason, if you are using Helm 2 on Kubernetes cluster older than 1.14.0, you will need to make tiller available on cluster using following steps:

```bash
# Skip these steps if you are using Helm v3 on Kubernetes cluster newer than v1.14.0
$ kubectl create serviceaccount tiller --namespace="kube-system"
$ kubectl create clusterrolebinding tiller --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
$ helm init --service-account=tiller
```

Now you can install the LogicMonitor Collectorset controller:

Get the configuration file downloaded from the LogicMonitor UI or you can create from the template [here](https://github.com/logicmonitor/k8s-helm-charts/blob/master/config-templates/Configuration.md#collectorset-controller).

Update configuration parameters in configuration file.

```bash
# Export the configuration file path & use it in th helm command.
$ export COLLECTORSET_CONTROLLER_CONF_FILE=<collectorset-controller-configuration-file-path>

$ helm upgrade \
  --install \
  --debug \
  --wait \
  --namespace="$NAMESPACE" \
  -f "$COLLECTORSET_CONTROLLER_CONF_FILE" \
  collectorset-controller logicmonitor/collectorset-controller
```

See the [configuration page](https://logicmonitor.github.io/k8s-argus/docs/configuration/) for a list of values the Collectorset Controller helm chart supports, and their
descriptions.

> Note: The Collectorset controller should be installed only once per cluster.

Next, install Argus:

Get the configuration file downloaded from the LogicMonitor UI or you can create from the template [here](https://github.com/logicmonitor/k8s-helm-charts/blob/master/config-templates/Configuration.md#argus).

Update configuration parameters in configuration file.

```bash
# Export the configuration file path & use it in th helm command.
$ export ARGUS_CONF_FILE=<argus-configuration-file-path>

$ helm upgrade \
  --install \
  --debug \
  --wait \
  --namespace="$NAMESPACE" \
  -f "$ARGUS_CONF_FILE" \
  argus logicmonitor/argus
```

See the [configuration page](https://logicmonitor.github.io/k8s-argus/docs/configuration/) for a list of values the Argus helm chart supports, and their descriptions.

> Note: Argus should be installed only once per cluster.

After installation is complete, you should make sure you have [DataSources](https://logicmonitor.github.io/k8s-argus/docs/monitoring/) in your account
that will start monitoring the resources in your cluster.

# Community

- To report bugs and/or submit feature requests, use [GitHub](https://github.com/logicmonitor/k8s-argus/issues).
