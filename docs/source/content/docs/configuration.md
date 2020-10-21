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

- **accessID (default: `""`):** The LogicMonitor API key ID.
- **accessKey (default: `""`):** The LogicMonitor API key.
- **account (default: `""`):** The LogicMonitor account name.
- **debug (default: `false`):** To enable verbose logging at debug level.

Optional Values:

- **enableRBAC (default: `true`):** Enable RBAC. If your cluster does not have
RBAC enabled, this value should be set to false.
- **etcdDiscoveryToken (default: `""`):** The public etcd discovery token used to add etcd hosts
 to the cluster device group.
 - **imagePullPolicy (default: `"Always"`):** The image pull policy of the Collectorset-controller container.
- **imageRepository (default: `"logicmonitor/collectorset-controller"`):** The image repository of the [Collectorset-controller](https://hub.docker.com/r/logicmonitor/collectorset-controller) container.
- **imageTag:** The image tag of the [Collectorset-controller](https://hub.docker.com/r/logicmonitor/collectorset-controller/tags) container.
- **proxyURL (default: `""`):** The Http/s proxy url.
- **proxyUser (default: `""`):** The Http/s proxy username.
- **proxyPass (default: `""`):** The Http/s proxy password.

Check the *[collectorset-controller-config.yaml](https://github.com/logicmonitor/k8s-helm-charts/blob/master/config-templates/Configuration.md#collectorset-controller)* for a complete list of values the Collectorset-controller helm chart supports.

# Configuring Argus via the Helm Chart

The Argus Helm chart supports the fololowing values:

Required Values:

- **accessID (default: `""`):** The LogicMonitor API key ID.
- **accessKey (default: `""`):** The LogicMonitor API key.
- **account (default: `""`):** The LogicMonitor account name.
- **clusterName (default: `""`):** A unique name given to the cluster's device group.
- **debug (default: `false`):** To enable verbose logging at debug level.
- **deleteDevices (default: `true`):** On a delete event, either delete from LogicMonitor or move the device to the `_delted` device group.
- **disableAlerting (default: `false`):** Disables LogicMonitor alerting for all the cluster resources.
- **collector.replicas (default: `1`):** The number of collectors to create and use with Argus.
- **collector.size (default: `""`):** The collector size to install. Can be nano, small, medium, or large.
- **collector.imageRepository (default: `logicmonitor/collector`):** The image repository of the [Collector](https://hub.docker.com/r/logicmonitor/collector) container.
- **collector.imageTag:** The image tag of the [Collector](https://hub.docker.com/r/logicmonitor/collector/tags) container.
- **collector.imagePullPolicy (default: `Always`):** The image pull policy of the Collector container.
- **collector.secretName (default: `"collector"`):** The Secret resource name of the collectors.

Optional Values:

- **debug (default: `false`):** To enable verbose logging at debug level.
- **deleteDevices (default: `true`):** On a delete event, either delete from
LogicMonitor or move the device to the `_deleted` device group.
- **disableAlerting (default: `false`):** Disables LogicMonitor alerting for all the cluster resources.
- **enableRBAC (default: `true`):** Enable RBAC. If your cluster does not have
RBAC enabled, this value should be set to false.
- **etcdDiscoveryToken:** The public etcd discovery token used to add etcd hosts
 to the cluster device group.
- **imagePullPolicy (default: `"Always"`):** The image pull policy of the Argus container.
- **imageRepository (default: `"logicmonitor/argus"`):** The image respository of the [Argus](https://hub.docker.com/r/logicmonitor/argus) container.
- **imageTag:** The image tag of the [Argus](https://hub.docker.com/r/logicmonitor/argus/tags) container.
- **proxyURL (default: `""`):** The Http/s proxy url.
- **proxyUser (default: `""`):** The Http/s proxy username.
- **proxyPass (default: `""`):** The Http/s proxy password.

Check the *[argus-config.yaml](https://github.com/logicmonitor/k8s-helm-charts/blob/master/config-templates/Configuration.md#argus)* for a complete list of values the Argus helm chart supports.

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

To configure the sensitive information, export the following environment
variables:

```bash
ACCESS_ID
ACCESS_KEY
ACCOUNT
ETCD_DISCOVERY_TOKEN
```
