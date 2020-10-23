---
title: "How It Works"
date: 2017-08-17T16:31:45-07:00
draft: false
menu:
  main:
    parent: Docs
    identifier: "How It Works"
    weight: 3
---

In this section we will dig into the lower level implementation of Argus to understand how it works, and provide those interested in contributing an introduction to the fundamentals of its design. An understanding of Go interfaces is recommended.

# Running Argus In-Cluster

Argus depends on communicating with the Kubernetes API Server. There are two ways to communicate with the API Server. In-cluster, and out-of-cluster. The `kubectl` CLI would be an example of out-of-cluster communication. Argus takes the former approach.

Running Argus in-cluster has advantages over running it out-of-cluster. For starters, you get all of the features that come with deploying an application in Kubernetes. Additionally, you can secure Argus using a `ServiceAccount` with `RBAC` policies allowing access only to what is required for Argus to function. Technically this _is_ possible out-of-cluster, but by being in-cluster, Kubernetes will take care of ensuring Argus has the `ServiceAccountToken` available at runtime. This simplifies things operationally and developmentally.

Finally, we need Argus on the same overlay network as the various Kubernetes resources. Since the collector comes with Argus, and the collector is on the overlay network, it can do its job without ever having to be Kubernetes aware.

# Watching Kubernetes Events

One of the basic functions of Argus is to represent the state of a Kubernetes cluster in LogicMonitor. To do that, it must be able to keep up with rapid changes of a constantly evolving cluster. Argus acheives this by registering event handlers for each resource we are instersted in representing in LogicMonitor. To understand how Argus can automate the management of various LogicMonitor resources, we need to understand what a `Controller` is. To quote the [documentation](https://kubernetes.io/docs/admin/kube-controller-manager/):

> In Kubernetes, a controller is a control loop that watches the shared state of the cluster through the apiserver and makes changes attempting to move the current state towards the desired state. Examples of controllers that ship with Kubernetes today are the replication controller, endpoints controller, namespace controller, and serviceaccounts controller.

The concept of a `Controller` is fundamental to Kubernetes and is at the core of its design. While Argus isn't a `Controller` in the sense that it _"makes changes [to the state of the cluster] attempting to move the current state towards the desired state"_, it _is_ a `Controller` in the sense that it moves a LogicMonitor account's state to match that of a cluster's state. Argus abstracts this into the notion of a `Watcher` that is responsible for watching Kubernetes events for a given resource and syncing the state to LogicMonitor.

# Implementing a Watcher

Now that we know about this event stream, let's look at what it takes to map resources in Kubernetes to objects in LogicMonitor. We start by first implenting the `Watcher` interface and then embedding a `Manager` in the concrete type implementing said interface. A `Watcher` is a simple interface that makes a concrete type compatible with the `NewInformer` function:

```
func (a *Argus) Watch() {
    getter := a.K8sClient.Core().RESTClient()
    for _, w := range a.Watchers {
        watchlist := cache.NewListWatchFromClient(getter, w.Resource(), v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(
            watchlist,
            w.ObjType(),
            time.Second*0,
            cache.ResourceEventHandlerFuncs{
                AddFunc:    w.AddFunc(),
                DeleteFunc: w.DeleteFunc(),
                UpdateFunc: w.UpdateFunc(),
            },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
    }
}
```

And we can see that the `Watcher` is defined as:

```
type Watcher interface {
    Resource() string // The Kubernetes resource that we want to watch (nodes, pods, services, etc.)
    ObjType() runtime.Object // A concrete type that is used for type assertion to an interface's underlying concrete value (Pod{}, Node{}, Service{}, etc.).
    AddFunc() func(obj interface{}) // A function that is responsible for handling add events for the given resource.
    UpdateFunc() func(oldObj, newObj interface{}) // A function that is responsible for handling update events for the given resource.
    DeleteFunc() func(obj interface{}) // A function that is responsible for handling delete events for the given resource.
}
```

With this simple function we can watch each Kubernetes resource we are interested in monitoring and provide custom logic for mapping it into LogicMonitor.

## The Manager

Now that we can watch events for a given resource, we need to implement the logic behind the add, update, and delete events. This is where we introduce the concept of a `Manager`. There are two functions of a `Manager`. First, a `Manager` must provide a way to build a LogicMonitor object given a Kubernetes resource object. Second, a `Manager` must ensure that the built object gets created in LogicMonitor. These concepts are abstracted into two interfaces, a `Builder` and a `Mapper`.

Let's imagine that we want to map a Kubernetes `Foo` resource into LogicMonitor as a `Bar` object.

```
type FooWatcher struct {
    BarManager
}

type BarManager interface {
    BarBuilder
    BarMapper
}
```

Here we can see that the Kuberentes `Foo` resource is mapped into LogicMonitor via a `Watcher` that implements the `BarManager` interface.
