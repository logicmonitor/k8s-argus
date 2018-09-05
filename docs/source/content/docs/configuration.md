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

# Configuring the Collecorset Controller via the Helm Chart

The Collectorset controller Helm chart supports the fololowing values:

Required Values:

- **accessID:** The LogicMonitor API key ID.
- **accessKey:** The LogicMonitor API key.
- **account:** The LogicMonitor account name.
- **clusterName:** A unique name given to the cluster's device group.

Optional Values:

- **enableRBAC (default: `true`):** Enable RBAC.
- **etcdDiscoveryToken:** The public etcd discovery token used to add etcd hosts to the cluster device group.
- **imagePullPolicy (default: `"Always"`):**
- **imageRepository (default: `"logicmonitor/argus"`):** The respository to use for the Argus docker image.
- **imageTag:** The collectorset-controller image tag to use.

# Configuring Argus via the Helm Chart

The Argus Helm chart supports the fololowing values:

Required Values:

- **accessID:** The LogicMonitor API key ID.
- **accessKey:** The LogicMonitor API key.
- **account:** The LogicMonitor account name.
- **clusterName:** A unique name given to the cluster's device group.
- **collector.replicas:** The number of collectors to create and use with Argus.
- **collector.size:** The collector size to install. Can be nano, small, medium, or large.

Optional Values:

- **debug (default: `false`):** Enable debug logging.
- **deleteDevices (default: `true`):** On a delete event, either delete from LogicMonitor or move the device to the `_delted` device group.
- **disableAlerting (default: `false`):** Disable alerting for all devices added.
- **enableRBAC (default: `true`):** Enable RBAC.
- **etcdDiscoveryToken:** The public etcd discovery token used to add etcd hosts to the cluster device group.
- **imagePullPolicy (default: `"Always"`):**
- **imageRepository (default: `"logicmonitor/argus"`):** The respository to use for the Argus docker image.
- **imageTag:** The argus container image tag to use.

# Configuring Argus Manually

In most applications there are generally two types of configuration options available. Options that do not contain sensitive information, and options that do contain sensitive information. Argus retrieves its' configuration from two different sources for each of these types.

 For non-sensitive configuration options, Arugs will read from a file on disk. For sensitive information, Arugs will read from environment variables.

To configure the non-sensitive information, create a YAML file located at `/etc/argus/config.yaml`. Here is an example file you can modify to your needs:

```yaml
cluster_name:
debug: false
delete_devices: true
disable_alerting: false
```

To configure the sensitive information, export the followin envionment variables:

```bash
ACCESS_ID
ACCESS_KEY
ACCOUNT
ETCD_DISCOVERY_TOKEN
```
