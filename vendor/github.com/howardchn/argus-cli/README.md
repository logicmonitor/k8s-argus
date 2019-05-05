# Argus Commandline Utility

This is a command utility for [argus project](https://logicmonitor.github.io/k8s-argus/) that help to setup kubernetes monitoring easily.

## Features
The original argus project is only allowing to install resources on your kubernetes environment and sync the resources into Santaba. It is complicated to uninstall it from your env and Santaba. Here are the details:

1. Device groups and resources within that reflect to your k8s resource structure
2. Collector groups and collectors based on the replicas number you set
3. Argus related resources (pods, crd, deployments, statefulset etc.)
4. CollectorSet-Controller related resources (pods, collectorsets etc.)

So this utility helps to clean those resources for you.

## Configurations



### Usage: argus-cli:
```
Usage:
  argus-cli [command]

Available Commands:
  help        Help about any command
  uninstall   uninstall argus related resources
  version     argus-cli version

Flags:
  -i, --accessId string    access id that is generated from santaba
  -k, --accessKey string   access key that is generated from santaba
  -h, --help               help for argus-cli
```

### Usage: argus-cli uninstall
```
Usage:
  argus-cli uninstall [flags]

Flags:
  -a, --account string    account name
  -c, --cluster string    cluster name
  -f, --confFile string   configure file (*.yaml)
  -h, --help              help for uninstall
  -m, --mode string       uninstall mode: [rest|helm|all], default: all (default "all")
  -g, --parentId int32    parent group id, default: 1 (default 1)

Global Flags:
  -i, --accessId string    access id that is generated from santaba
  -k, --accessKey string   access key that is generated from santaba
```

Example:
```
$ ./argus-cli uninstall --accessId="[Access ID]" --accessKey="[Access Key]" --clusterName="[Cluster Name]" --account="[Company Name]" --parentId=[Parent Group Id]
```

### Usage: argus-cli uninstall -f
Example:
```
$ ./argus-cli uninstall -f ./conf.yaml
```

Configuration file format:
```yaml
accessId:   [replace with accessId          *required]
accessKey:  [replace with accessKey         *required]
account:    [replace with account           *required]
cluster:    [replace with cluster name      *required]
parentId:   [replace with parent group id]
mode:       [replace with mode (all|rest|helm)]
```

## Build Executable
Feel free to clone this project and build for your specific OS.

### Precondition:
1. Setup the go environment on your machine. Refer the [installation instruction](https://golang.org/doc/install).
2. Kubernetes environment setup.
3. Helm installed.

### Build Detail
Use command to build for various platform

Platform | Build Command
---------|--------------
macOS | CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o argus-cli main.go
Linux | CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o argus-cli main.go
Win   | CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o argus-cli main.go

Use makefile to build for various platform
```bash
# cd to the project
make darwin
make linux
make win
```

## Contact
* [Email](mailto:howardch@outlook.com)
* [Linkedin](https://www.linkedin.com/in/howard-chen-328493142/)
