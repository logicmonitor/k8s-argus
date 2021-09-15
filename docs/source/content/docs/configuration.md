---
title: Configuration
date: '2017-08-17T00:54:55.000Z'
draft: false
menu:
  main:
    parent: Docs
    identifier: Configuration
    weight: 1
---

# configuration

## Configuring the Collectorset Controller via the Helm Chart

The Collectorset controller Helm chart supports the following values:

Required Values:

- **accessID (default: `""`):** The LogicMonitor API key ID.
- **accessKey (default: `""`):** The LogicMonitor API key.
- **account (default: `""`):** The LogicMonitor account name.
- **debug (default: `false`):** To enable verbose logging at debug level.

Optional Values:

- **enableRBAC (default: `true`):** Enable RBAC. If your cluster does not have RBAC enabled, this value should be set to false.
- **etcdDiscoveryToken (default: `""`):** The public etcd discovery token used to add etcd hosts to the cluster device group.
- **imagePullPolicy (default: `"Always"`):** The image pull policy of the Collectorset-controller container.
- **imageRepository (default: `"logicmonitor/collectorset-controller"`):** The image repository of the [Collectorset-controller](https://hub.docker.com/r/logicmonitor/collectorset-controller) container.
- **imageTag:** The image tag of the [Collectorset-controller](https://hub.docker.com/r/logicmonitor/collectorset-controller/tags) container.
- **proxyURL (default: `""`):** The Http/s proxy url.
- **proxyUser (default: `""`):** The Http/s proxy username.
- **proxyPass (default: `""`):** The Http/s proxy password.
- **nodeSelector (default: `{}`):** It provides the simplest way to run Pod on particular Node(s) based on labels on the node.
- **affinity (default: `{}`):** It allows you to constrain which nodes your pod is eligible to be scheduled on.
- **priorityClassName (default: `""`):** The priority class name for Pod priority. If this parameter is set then user must have PriorityClass resource created otherwise Pod will be rejected.
- **tolerations (default: `[]`):** Tolerations are applied to pods, and allow the pods to schedule onto nodes with matching taints.

* **imageTag:** The collectorset-controller \[image tag\] \([https://hub.docker.com/r/logicmonitor/collectorset-controller/tags/](https://hub.docker.com/r/logicmonitor/collectorset-controller/tags/)\) to use.
* **proxyURL \(default: `""`\):** The Http/s proxy url.
* **proxyUser \(default: `""`\):** The Http/s proxy username.
* **proxyPass \(default: `""`\):** The Http/s proxy password.

## Configuring Argus via the Helm Chart

The Argus Helm chart supports the fololowing values:

Required Values:

- **accessID (default: `""`):** The LogicMonitor API key ID.
- **accessKey (default: `""`):** The LogicMonitor API key.
- **account (default: `""`):** The LogicMonitor account name.
- **clusterName (default: `""`):** A unique name given to the cluster's device group.
- **debug (default: `false`):** To enable verbose logging at debug level.
- **deleteDevices (default: `true`):** On a delete event, either delete from LogicMonitor or move the device to the `_deleted` device group.
- **disableAlerting (default: `false`):** Disables LogicMonitor alerting for all the cluster resources.
- **collector.replicas (default: `1`):** The number of collectors to create and use with Argus.
- **collector.size (default: `""`):** The collector size to install. Can be nano, small, medium, or large.
- **collector.imageRepository (default: `logicmonitor/collector`):** The image repository of the [Collector](https://hub.docker.com/r/logicmonitor/collector) container.
- **collector.imageTag:** The image tag of the [Collector](https://hub.docker.com/r/logicmonitor/collector/tags) container.
- **collector.imagePullPolicy (default: `Always`):** The image pull policy of the Collector container.
- **collector.secretName (default: `"collector"`):** The Secret resource name of the collectors.

Optional Values:

- **enableRBAC (default: `true`):** Enable RBAC. If your cluster does not have RBAC enabled, this value should be set to false.
- **clusterGroupID (default: `0`):** A parent group id of the cluster's device group.
- **etcdDiscoveryToken (default: `""`):** The public etcd discovery token used to add etcd hosts to the cluster device group.
- **imagePullPolicy (default: `"Always"`):** The image pull policy of the Argus container.
- **imageRepository (default: `"logicmonitor/argus"`):** The image respository of the [Argus](https://hub.docker.com/r/logicmonitor/argus) container.
- **imageTag:** The image tag of the [Argus](https://hub.docker.com/r/logicmonitor/argus/tags) container.
- **proxyURL (default: `""`):** The Http/s proxy url.
- **proxyUser (default: `""`):** The Http/s proxy username.
- **proxyPass (default: `""`):** The Http/s proxy password.
- **nodeSelector (default: `{}`):** It provides the simplest way to run Pod on particular Node(s) based on labels on the node.
- **affinity (default: `{}`):** It allows you to constrain which nodes your pod is eligible to be scheduled on.
- **priorityClassName (default: `""`):** The priority class name for Pod priority. If this parameter is set then user must have PriorityClass resource created otherwise Pod will be rejected.
- **tolerations (default: `[]`):** Tolerations are applied to pods, and allow the pods to schedule onto nodes with matching taints.
- **filters.pod (default: `""`):** The filtered expression for Pod device type. Based on this parameter, Pods would be added/deleted for discovery on LM.
- **filters.service (default: `""`):** The filtered expression for Service device type. Based on this parameter, Services would be added/deleted for discovery on LM.
- **filters.node (default: `""`):** The filtered expression for Node device type. Based on this parameter, Nodes would be added/deleted for discovery on LM.
- **filters.deployment (default: `""`):** The filtered expression for Deployment device type. Based on this parameter, Deployments would be added/deleted for discovery on LM.
- **collector.groupID (default: `0`):** The ID of the group of the collectors.
- **collector.escalationChainID (default: `0`):** The ID of the escalation chain of the collectors.
- **collector.collectorVersion (default: `0`):** The version of the collectors.
- **collector.useEA (default: `false`):** On a collector downloading event, either download the latest EA version or the latest GD version.
- **collector.proxyURL (default: `""`):** The Http/s proxy url of the collectors.
- **collector.proxyUser (default: `""`):** The Http/s proxy username of the collectors.
- **collector.proxyPass (default: `""`):** The Http/s proxy password of the collectors.
- **collector.priorityClassName (default: `""`):** The priority class name for Pod priority of the collector. If this parameter is set then user must have PriorityClass resource created otherwise Pod will be rejected.
- **collector.tolerations (default: `[]`):** Tolerations are applied to pods, and allow the pods to schedule onto nodes with matching taints.

# Configuring Argus Manually

In most applications there are generally two types of configuration options
available. Options that do not contain sensitive information, and options that
do contain sensitive information. Argus retrieves its' configuration from two
different sources for each of these types.

 For non-sensitive configuration options, Argus will read from a file on disk.
 For sensitive information, Argus will read from environment variables.

To configure the non-sensitive information, create a YAML file located at
`/etc/argus/config.yaml`. Here is an example file you can modify to your needs:

```yaml
cluster_group_id:
cluster_name:
debug: false
delete_devices: true
disable_alerting: false
```

To configure the sensitive information, export the following environment variables:

```bash
ACCESS_ID
ACCESS_KEY
ACCOUNT
ETCD_DISCOVERY_TOKEN
```

