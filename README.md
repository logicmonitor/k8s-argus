<p align="center">
  <h1 align="center">Argus</h1>
  <p align="center">Automated Kubernetes monitoring.</p>
  <p align="center"><sub> <i>Powered by LogicMonitor</i></sub></p>
  <p align="center">
    <a href="https://godoc.org/github.com/logicmonitor/argus"><img alt="GoDoc" src="http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
  </p>
  <p align="center">
    <a href="https://github.com/logicmonitor/argus/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/logicmonitor/argus.svg?style=flat-square"></a>
    <a href="https://github.com/logicmonitor/argus/releases/latest"><img alt="GitHub (pre-)release" src="https://img.shields.io/github/release/logicmonitor/argus/all.svg?style=flat-square"></a>
  </p>
</p>

---

**Argus** is a tool for monitoring of Kubernetes with [LogicMonitor](https://www.logicmonitor.com). Some of the key features of Argus are:
-   **Automated discovery:** Discovers nodes, services, and pods, and adds them to your LogicMonitor account automatically.
-   **Datasources applied to Kubernetes labels:** Resources are added along with all of its' associated tags, allowing for custom datasources to be applied based on Kubernetes labels.
-   **Real-time** By leveraging the Kubernetes API event stream, Argus keeps the monitoring of your cluster up-to-date and accurate.

Getting Started
---------------
Argus is available as a Helm chart. Please see the installation instructions in the [LogicMonitor charts](https://github.com/logicmonitor/k8s-charts) repository.
> **Note:** Argus opensource software and currently in alpha phase. The LogicMonitor support team will not assist in any issues related to Argus.

Developing Argus
----------------
To compile the binary and build the Docker image, run:
```
$ make
```

If you would like to run the same tests used in the CI build, run:
```
$ make test
```
> **Note:** Running tests before submitting a PR is highly recommended.

### License
[![license](https://img.shields.io/github/license/logicmonitor/k8s-argus.svg?style=flat-square)](https://github.com/logicmonitor/k8s-argus/blob/master/LICENSE)
