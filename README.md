<p align="center"><img src="./logo.png"></p>

---

<p align="center">
  <a href="https://gitter.im/k8s-argus/Lobby"><img alt="GoDoc" src="https://img.shields.io/gitter/room/k8s-argus/Lobby.svg?style=flat-square"></a>
  <a href="https://godoc.org/github.com/logicmonitor/k8s-argus"><img alt="GoDoc" src="http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
  <a href="https://travis-ci.org/logicmonitor/k8s-argus"><img alt="Travis" src="https://img.shields.io/travis/logicmonitor/k8s-argus.svg?style=flat-square"></a>
  <a href="https://codecov.io/gh/logicmonitor/k8s-argus"><img alt="Codecov" src="https://img.shields.io/codecov/c/github/logicmonitor/k8s-argus.svg?style=flat-square"></a>
  <a href="https://goreportcard.com/report/github.com/logicmonitor/k8s-argus"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/logicmonitor/k8s-argus?style=flat-square"></a>
  <a href="https://github.com/logicmonitor/k8s-argus/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/logicmonitor/argus.svg?style=flat-square"></a>
  <a href="https://github.com/logicmonitor/k8s-argus/releases/latest"><img alt="GitHub (pre-)release" src="https://img.shields.io/github/release/logicmonitor/argus/all.svg?style=flat-square"></a>
</p>

**Argus** is a tool for monitoring Kubernetes with
[LogicMonitor](https://www.logicmonitor.com). Some of the key features of Argus are:

- **Automated Resource Discovery:** Leverages Kubernetes events and LogicMonitor's API to provide real-time accuracy of
  a cluster's resources in LogicMonitor. Discovers etcd members, cluster Nodes, Services, and Pods, and automates the
  management of their lifecycle as
  [Resources](https://www.logicmonitor.com/support/devices/) in LogicMonitor.
- **Comprehensive Monitoring:** Dockerized LogicMonitor Collectors running in a Stateful Set and managed by Argus
  collect data via the Kubernetes API for Nodes, Pods, Services, and Containers. Additionally, you can leverage
  LogicMonitor DataSources to monitor your applications running within the cluster.

See the [documentation](https://logicmonitor.github.io/k8s-argus) to discover more about Argus.

Developing Argus
----------------
To build Argus, run:

```
$ make
```

> **Note:** The Dockerfile uses multi-stage builds. Docker 17.05.0 or greater is required.

To build the documentation, run:

```
$ make docs
```

> **Note:** [Hugo](https://github.com/gohugoio/hugo) is required to generate the documentation.



### License

[![license](https://img.shields.io/github/license/logicmonitor/k8s-argus.svg?style=flat-square)](https://github.com/logicmonitor/k8s-argus/blob/master/LICENSE)
