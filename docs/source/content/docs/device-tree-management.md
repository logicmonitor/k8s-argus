---
title: "Device Tree Management"
date: 2017-08-16T17:55:12-07:00
draft: false
menu:
  main:
    parent: Docs
    identifier: "Device Tree Management"
    weight: 2
---

One of the main features of Argus is its ability to represent a Kubernetes cluster in LogicMonitor, and to keep that representation up to date and accurate. Argus achieves this by managing the following.

# Device Groups

Argus has an opinionated way of representing a cluster in the LogicMonitor Device Tree.
It will first create a top level device group with the name of your cluster as specified in the arguments to the chart.
Under this top level device group, additional device groups will be created for each resource type that Argus supports. To name but a few:

- **Nodes:** The nodes that the cluster is comprised of.
- **Pods:** Pods running in the cluster.
- **Services:** Services running in the cluster.

For namespaced resources (Pods, Services), a device group will be created for each namespace and placed inside the respective resource device group.
For clusters that use an external etcd cluster, an **Etcd** device group will be created in the cluster's root device group.

# Devices

Argus subscribes to the Kubernetes API event stream for each resource type mentioned in [Device Groups](#device-groups). Upon receiving an event, Argus decides on how to represent the event in LogicMonitor. In the case of nodes, pods, and services, Argus will create a device.

## Device Properties

Argus will map Kubernetes metadata to device properties.
This means that [Datasources](https://www.logicmonitor.com/support/datasources/) can be applied based on things like [labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).
