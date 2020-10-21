---
title: "Getting Started"
date: 2017-08-12T16:20:39-07:00
draft: false
---

The simplest way to get started with Argus is to install it using [Helm](https://github.com/kubernetes/helm). Prior to installation, you will need a ClusterRoleBinding & ServiceAccount for tiller:

```bash
$ kubectl create serviceaccount tiller --namespace="kube-system"
$ kubectl create clusterrolebinding tiller --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
$ helm init --service-account=tiller
```

> Note: You can skip above steps if you are installing helm charts with Helm v3.

You'll also need to add the LogicMonitor chart repository:

```bash
$ helm repo add logicmonitor https://logicmonitor.github.com/k8s-helm-charts
```

Now you can install the LogicMonitor Collectorset controller:
Create *[collectorset-controller-configuration.yaml](https://github.com/logicmonitor/k8s-helm-charts/blob/master/config-templates/Configuration.md#collectorset-controller)* file and add required values in it then pass the file path in the helm command.

```bash
$ helm upgrade \
  --install \
  --debug \
  --wait \
  --namespace="$NAMESPACE" \
  -f collectorset-controller-configuration.yaml \
  collectorset-controller logicmonitor/collectorset-controller
```

See the [configuration page](https://logicmonitor.github.io/k8s-argus/docs/configuration/) for a list of values the Collectorset Controller helm chart supports, and their
descriptions.

> Note: The Collectorset controller should be installed only once per cluster.

Next, install Argus:
Create *[argus-configuration.yaml](https://github.com/logicmonitor/k8s-helm-charts/blob/master/config-templates/Configuration.md#argus)* file and add required values in it then pass the file path in the helm command.

```bash
$ helm upgrade \
  --install \
  --debug \
  --wait \
  --namespace="$NAMESPACE" \
  -f argus-configuration.yaml \
  argus logicmonitor/argus
```

See the [configuration page](https://logicmonitor.github.io/k8s-argus/docs/configuration/) for a list of values the Argus helm chart supports, and their descriptions.

> Note: Argus should be installed only once per cluster.

After installation is complete, you should make sure you have [DataSources](https://logicmonitor.github.io/k8s-argus/docs/monitoring/) in your account
that will start monitoring the resources in your cluster.

# Community

- To report bugs and/or submit feature requests, use [GitHub](https://github.com/logicmonitor/k8s-argus/issues).
