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

* **accessID:** The LogicMonitor API key ID.
* **accessKey:** The LogicMonitor API key.
* **account:** The LogicMonitor account name.
* **clusterName:** A unique name given to the cluster's device group.

Optional Values:

* **enableRBAC \(default: `true`\):** Enable RBAC. If your cluster does not have

  RBAC enabled, this value should be set to false.

* **etcdDiscoveryToken:** The public etcd discovery token used to add etcd hosts

  to the cluster device group.

* **imagePullPolicy \(default: `"Always"`\):**
* **imageRepository \(default: `"logicmonitor/collectorset-controller"`\):** The

  respository to use for the collectorset-controller docker image.

* **imageTag:** The collectorset-controller \[image tag\] \([https://hub.docker.com/r/logicmonitor/collectorset-controller/tags/](https://hub.docker.com/r/logicmonitor/collectorset-controller/tags/)\) to use.
* **proxyURL \(default: `""`\):** The Http/s proxy url.
* **proxyUser \(default: `""`\):** The Http/s proxy username.
* **proxyPass \(default: `""`\):** The Http/s proxy password.

## Configuring Argus via the Helm Chart

The Argus Helm chart supports the fololowing values:

Required Values:

* **accessID:** The LogicMonitor API key ID.
* **accessKey:** The LogicMonitor API key.
* **account:** The LogicMonitor account name.
* **clusterGroupID:** A parent group id of the cluster's device group.
* **clusterName:** A unique name given to the cluster's device group.
* **collector.replicas:** The number of collectors to create and use with Argus.
* **collector.size:** The collector size to install. Can be nano, small, medium,

  or large.

Optional Values:

* **debug \(default: `false`\):** Enable debug logging.
* **deleteDevices \(default: `true`\):** On a delete event, either delete from

  LogicMonitor or move the device to the `_delted` device group.

* **disableAlerting \(default: `false`\):** Disable alerting for all devices added.
* **enableRBAC \(default: `true`\):** Enable RBAC. If your cluster does not have

  RBAC enabled, this value should be set to false.

* **etcdDiscoveryToken:** The public etcd discovery token used to add etcd hosts

  to the cluster device group.

* **imagePullPolicy \(default: `"Always"`\):**
* **imageRepository \(default: `"logicmonitor/argus"`\):** The respository to use

  for the Argus docker image.

* **imageTag:** The argus container \[image tag\] \([https://hub.docker.com/r/logicmonitor/argus/tags/](https://hub.docker.com/r/logicmonitor/argus/tags/)\) to use.
* **proxyURL \(default: `""`\):** The Http/s proxy url.
* **proxyUser \(default: `""`\):** The Http/s proxy username.
* **proxyPass \(default: `""`\):** The Http/s proxy password.

## Configuring Argus Manually

In most applications there are generally two types of configuration options available. Options that do not contain sensitive information, and options that do contain sensitive information. Argus retrieves its' configuration from two different sources for each of these types.

For non-sensitive configuration options, Argus will read from a file on disk. For sensitive information, Argus will read from environment variables.

To configure the non-sensitive information, create a YAML file located at `/etc/argus/config.yaml`. Here is an example file you can modify to your needs:

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

