---
title: Monitoring
date: '2017-08-12T23:20:39.000Z'
draft: false
menu:
  main:
    parent: Docs
    identifier: Monitoring
    weight: 1
---

Argus adds Kubernetes resources into LogicMonitor, but the DataSources that
apply to those resources are responsible for monitoring. LogicMonitor has a set
 of developed modules for monitoring Kubernetes, which you can import via the
 [LM Exchange](https://www.logicmonitor.com/support/settings/logicmodules/lm-exchange/)
 (Settings | DataSources | Add | From LogicMonitor Exchange). Filter Items with word “Kubernetes”. Add all logicmodules as per your requirements.

If kubernetes cluster is configured with Istio then you can add “Istio kubernetes” logicmodules as well.

* Kubernetes\_Nodes: PR4F33
* Kubernetes\_Node: 267H63
* Kubernetes\_Healthz: KZ463J
* Kubernetes\_ControlPlane: N3GZNX
* Kubernetes\_Service: HPJPRT
* Kubernetes\_Scheduler: FCPJNH
* Kubernetes\_Pod: P9TT2W
* Kubernetes\_Container: 3AAJZX
* Kubernetes\_PingK8s: 4N99FE

Beyond the health and performance of your Kubernetes Cluster resources, you can
 use LogicMonitor DataSources to monitor your applications running in
 Kubernetes. If it's a standard application, LogicMonitor likely already has a
 DataSource that can be used. For custom applications, you should
 [create your own DataSource](https://www.logicmonitor.com/support/datasources/creating-managing-datasources
  ) that is specific to your application. Either way, there are a few
  requirements to ensure that these DataSources will work well for applications
   running in Kubernetes.

1. **Applies To:**
  DataSources usually apply based on property values, such as a *system.categories* or *system.sysinfo* value. 
  For applications running in Kubernetes, the best option is to have Datasources apply based on labels.
  Argus adds labels as resource properties (prepended with kubernetes.label.), and you can reference these labels in the Applies To field for any DataSources. 
  For example, you may use a label app=shoppingcart, and have a DataSource that monitors shopping cart performance with an Applies To of *kubernetes.label.app=="shoppingcart"*.

2. **Referencing hostname and IP addresses:**
  Most LogicMonitor DataSources use ##HOSTNAME## to reference the IP, DNS, or system name for data collection. 
  This leverages the *system.hostname* property value in LogicMonitor. 
  Argus sets both *system.hostname* and *system.ips* to the Pod IP except hostNetwork enabled Pods. Argus sets *system.hostname* to the Pod name for hostNetwork enabled Pods. 
  As such, DataSources that monitor applications running in Kubernetes should use ##system.ips## instead of ##HOSTNAME##.
