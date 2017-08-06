<p align="center">
  <p align="center">
    <a href=""><img alt="GoDoc" src="./logo.png"></a>

  </p>
  <p align="center">
    <a href="https://godoc.org/github.com/logicmonitor/k8s-argus"><img alt="GoDoc" src="http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>

  </p>
  <p align="center">
  <a href="https://travis-ci.org/logicmonitor/k8s-argus"><img alt="Travis" src="https://img.shields.io/travis/logicmonitor/k8s-argus.svg?style=flat-square"></a>
  <a href="https://codecov.io/gh/logicmonitor/k8s-argus"><img alt="Codecov" src="https://img.shields.io/codecov/c/github/logicmonitor/k8s-argus.svg?style=flat-square"></a>
  <a href="https://goreportcard.com/report/github.com/logicmonitor/k8s-argus"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/logicmonitor/k8s-argus?style=flat-square"></a>
</p>
  <p align="center">
    <a href="https://github.com/logicmonitor/k8s-argus/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/logicmonitor/argus.svg?style=flat-square"></a>
    <a href="https://github.com/logicmonitor/k8s-argus/releases/latest"><img alt="GitHub (pre-)release" src="https://img.shields.io/github/release/logicmonitor/argus/all.svg?style=flat-square"></a>
  </p>
</p>

---

**Argus** is a tool for monitoring Kubernetes with [LogicMonitor](https://www.logicmonitor.com). Some of the key features of Argus are:
-   **Automated discovery:** Discovers nodes, services, and pods, and adds them as devices to your LogicMonitor account automatically.
-   **Real-time:** By leveraging the Kubernetes API event stream, Argus keeps everything a collector needs to know about your cluster up-to-date and accurate.
-   **Granular monitoring:** By adding resources as devices, data sources can use an IP address to talk directly to them. And since all Kubernetes labels are added as `system.categories` device properties, data sources can be applied based on lables.
-   **Organized visual represention:** Clusters are represented in the device tree as a collection of dynamic [Device Groups](https://www.logicmonitor.com/support/devices/device-groups/device-groups-overview/) and organizes the devices by resource type and namespace.

Getting Started
---------------
Argus is available as a Helm chart. Please see the installation instructions in the [LogicMonitor charts](https://github.com/logicmonitor/k8s-helm-charts/tree/master/argus) repository.
> **Note:** Argus opensource software and currently in alpha phase. The LogicMonitor support team will not assist in any issues related to Argus.

Developing Argus
----------------
To build Docker image, run:
```
$ make
```
> **Note:** The Dockerfile uses multi-stage builds. Docker 17.05.0 or greater is required.

### License
[![license](https://img.shields.io/github/license/logicmonitor/k8s-argus.svg?style=flat-square)](https://github.com/logicmonitor/k8s-argus/blob/master/LICENSE)
