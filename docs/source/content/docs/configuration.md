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

# Configuring Argus via the Helm Chart

The Helm chart supports the fololowing values:

Required Values:

- **accessID:** The LogicMonitor API key ID.
- **accessKey:** The LogicMonitor API key.
- **account:** The LogicMonitor account name.
- **clusterName:** A unique name given to the cluster's device group.
- **collectorDescription:** A unique collector description used to look up a collector dynamically.
- **collectorVersion:** The collector version to install.

Optional Values:

- **collectorEscalationChainID:** The ID of the escalation chain to use for collector down alerts.
- **collectorImageTag (default: `"latest"`):** The collector image tag.
- **collectorSize (default: `"small"`):** The collector size to install. Can be nano, small, medium, or large.
- **debug (default: `false`):** Enable debug logging.
- **deleteDevices (default: `true`):** On a delete event, either delete from LogicMonitor or move the device to the `_delted` device group.
- **disableAlerting (default: `false`):** Disable alerting for all devices added.
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
collector_description:
collector_escalation_chain_id:
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
