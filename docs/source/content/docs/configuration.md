---
title: "Configuration"
date: 2017-08-16T17:54:55-07:00
draft: false
menu:
  main:
    parent: Docs
    identifier: "Configuration"
    weight: 1
---

# Configuring the Collectorset Controller via the Helm Chart

The Collectorset controller Helm chart supports the following values:

Required Values:

* **accessID (default: `""`):** The LogicMonitor API key ID.
* **accessKey (default: `""`):** The LogicMonitor API key.
* **account (default: `""`):** The LogicMonitor account name.
* **debug (default: `false`):** To enable verbose logging at debug level.

Optional Values:

* **enableRBAC (default: `true`):** Enable RBAC. If your cluster does not have RBAC enabled, this value should be set to false.
* **etcdDiscoveryToken (default: `""`):** The public etcd discovery token used to add etcd hosts to the cluster device group.
* **imageRepository (default: `"logicmonitor/collectorset-controller"`):** The image repository of the [Collectorset-controller](https://hub.docker.com/r/logicmonitor/collectorset-controller) container.
* **imageTag:** The image tag of the [Collectorset-controller](https://hub.docker.com/r/logicmonitor/collectorset-controller/tags) container.
* **imagePullPolicy (default: `"Always"`):** The image pull policy of the Collectorset-controller container.
* **proxyURL (default: `""`):** The Http/s proxy url.
* **proxyUser (default: `""`):** The Http/s proxy username.
* **proxyPass (default: `""`):** The Http/s proxy password.
* **nodeSelector (default: `{}`):** It provides the simplest way to run Pod on particular Node(s) based on labels on the node.
* **affinity (default: `{}`):** It allows you to constrain which nodes your pod is eligible to be scheduled on.
* **priorityClassName (default: `""`):** The priority class name for Pod priority. If this parameter is set then user must have PriorityClass resource created otherwise Pod will be rejected.
* **tolerations (default: `[]`):** Tolerations are applied to pods, and allow the pods to schedule onto nodes with matching taints.
* **labels (default: `{}`):** Labels to apply on all objects created by Collectorset controller.
* **annotations (default: `{}`):** Annotations to apply on all objects created by Collectorset controller.
* **ignore_ssl (default: `false`):** Set flag to ignore ssl/tls validation.

# Configuring Argus via the Helm Chart

The Argus Helm chart supports the following values:

Required Values:

* **accessID (default: `""`):** The LogicMonitor API key ID.
* **accessKey (default: `""`):** The LogicMonitor API key.
* **account (default: `""`):** The LogicMonitor account name.
* **clusterName (default: `""`):** A unique name given to the cluster's device group.
* **logLevel (default: `"info"`):** Set Argus Log Level.
* **deleteDevices (default: `true`):** On a delete event, either delete from LogicMonitor or move the device to the `_deleted` device group.
* **disableAlerting (default: `false`):** Disables LogicMonitor alerting for all the cluster resources.
* **collector.replicas (default: `1`):** The number of collectors to create and use with Argus.
* **collector.size (default: `""`):** The collector size to install. Can be nano, small, medium, or large.
* **collector.imageRepository (default: `logicmonitor/collector`):** The image repository of the [Collector](https://hub.docker.com/r/logicmonitor/collector) container.
* **collector.imageTag:** The image tag of the [Collector](https://hub.docker.com/r/logicmonitor/collector/tags) container.
* **collector.imagePullPolicy (default: `Always`):** The image pull policy of the Collector container.
* **collector.secretName (default: `"collector"`):** The Secret resource name of the collectors.

Optional Values:

* **enableRBAC (default: `true`):** Enable RBAC. If your cluster does not have RBAC enabled, this value should be set to false.
* **clusterGroupID (default: `0`):** A parent group id of the cluster's device group.
* **etcdDiscoveryToken (default: `""`):** The public etcd discovery token used to add etcd hosts to the cluster device group.
* **imageRepository (default: `"logicmonitor/argus"`):** The image respository of the [Argus](https://hub.docker.com/r/logicmonitor/argus) container.
* **imageTag:** The image tag of the [Argus](https://hub.docker.com/r/logicmonitor/argus/tags) container.
* **imagePullPolicy (default: `"Always"`):** The image pull policy of the Argus container.
* **proxyURL (default: `""`):** The Http/s proxy url.
* **proxyUser (default: `""`):** The Http/s proxy username.
* **proxyPass (default: `""`):** The Http/s proxy password.
* **nodeSelector (default: `{}`):** It provides the simplest way to run Pod on particular Node(s) based on labels on the node.
* **affinity (default: `{}`):** It allows you to constrain which nodes your pod is eligible to be scheduled on.
* **priorityClassName (default: `""`):** The priority class name for Pod priority. If this parameter is set then user must have PriorityClass resource created otherwise Pod will be rejected.
* **tolerations (default: `[]`):** Tolerations are applied to pods, and allow the pods to schedule onto nodes with matching taints.
* **labels (default: `{}`):** Labels to apply on all objects created by Argus.
* **annotations (default: `{}`):** Annotations to apply on all objects created by Argus.
* **ignore_ssl (default: `false`):** Set flag to ignore ssl/tls validation.
* **enableNewResourceTree (default: `false`):** Flag to enable new resource tree to put all resources of a namespace in a single resource group.
* **enableNamespacesDeletedGroups (default: `false`):** Flag is used when #enableNewResourceTree is true - to create _deleted group in its individual namespace groups when scheduled delete is enabled.
* **registerGenericFilter (default: `false`):**  Flag to register generic filter based on resource label => "logicmonitor/monitoring": "disable".
* **app_intervals.periodic_sync_interval (default: `30m`):** AppIntervals defines time intervals for periodic sync, periodic delete and in memory cache resync operations.
* **app_intervals.periodic_delete_interval (default: `10m`):**  AppIntervals defines time intervals for periodic sync, periodic delete and in memory cache resync operations.
* **app_intervals.cache_sync_interval (default: `1h`):**  AppIntervals defines time intervals for periodic sync, periodic delete and in memory cache resync operations.
* **device_group_props.cluster:** Device group properties for cluster.
* **device_group_props.pods:** Device group properties for pods.
* **device_group_props.services (default: `[]`):** Contains device group properties for services.
* **device_group_props.deployments (default: `[]`):** Contains device group properties for deployments.
* **device_group_props.nodes (default: `[]`):** Contains device group properties for nodes.
* **device_group_props.etcd (default: `[]`):** Contains device group properties for etcd.
* **device_group_props.hpas (default: `[]`):** Contains device group properties for HorizontalPodAutoscalers.
* **filters.pod (default: `[]`):** The filtered expression for Pod device type. Based on this parameter, Pods would be added/deleted for discovery on LM.
* **filters.service (default: `[]`):** The filtered expression for Service device type. Based on this parameter, Services would be added/deleted for discovery on LM.
* **filters.node (default: `[]`):** The filtered expression for Node device type. Based on this parameter, Nodes would be added/deleted for discovery on LM.
* **filters.deployment (default: `[]`):** The filtered expression for Deployment device type. Based on this parameter, Deployments would be added/deleted for discovery on LM.
* **filters.hpa (default: `[]`):** The filtered expression for HorizontalPodAutoscaler resource type. Based on this parameter,HorizontalPodAutoscalers would be added/deleted for discovery on LM.
* **openmetrics.port (default: `2112`):** Openmetrics config for Argus metrics collection.
* **collector.groupID (default: `0`):** The ID of the group of the collectors.
* **collector.escalationChainID (default: `0`):** The ID of the escalation chain of the collectors.
* **collector.collectorVersion (default: `0`):** The version of the collectors.
* **collector.useEA (default: `false`):** On a collector downloading event, either download the latest EA version or the latest GD version.
* **collector.proxyURL (default: `""`):** The Http/s proxy url of the collectors.
* **collector.proxyUser (default: `""`):** The Http/s proxy username of the collectors.
* **collector.proxyPass (default: `""`):** The Http/s proxy password of the collectors.
* **collector.annotations (default: `{}`):** annotations to add on collector statefulset.
* **collector.labels (default: `{}`):** labels to add on collector statefulset.
* **collector.statefulsetspec:** Holds statefulset specification template which contains nodeSelector, tolerations and priorityClassName.
* **disableResourceMonitoring:** List of resources to disable monitoring.
* **disableResourceAlerting:** List of resources to disable alerting.
* **replicas (default: `1`):** Argus replicas.

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
clusterGroupID:
clusterName:
deleteDevices: true
disableAlerting: false
```

To configure the sensitive information, export the following environment
variables:

```bash
ACCESS_ID
ACCESS_KEY
ACCOUNT
ETCD_DISCOVERY_TOKEN
```
