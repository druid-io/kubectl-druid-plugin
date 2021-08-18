## kubectl-druid-plugin
- Kubectl plugin to simplify operations on apache druid CR.

## Prerequisite
- Apache Druid CRD defination to be installed.
- An apache druid CR managed apache druid cluster running. 
- https://github.com/druid-io/druid-operator/tree/master/deploy/crds

## Getting Started
- NOTE: go version 1.15
- ```go build -o kubectl-druid```
- ```mv kubectl-druid /usr/local/bin```
- ```kubectl druid --help```

#### Build using Docker

- ```make build```
- ```mv build/pkg/kubectl-druid-${OS}-${ARCH} /usr/local/bin/kubectl-druid```
- ```kubectl druid --help```

#### Download the artifacts

- Go to Github Actions
- Download the artifact zip file.

## Commands

- List All Druid CR's in a k8s cluster
```
kubectl druid list
```

- List Druid CR's in a namespace
```
kubectl druid list --namespace <namespace>
```

- Get Druid Nodes's in a namespace for a specific cr
```
kubectl druid get nodes --cr <cr>--namespace <namespace>
```

- Scale Druid Replicas for a specific druid cr node in a namespace
```
kubectl druid scale --cr <cr> --namespace <namespace> --node middlemanager --replicas 4
```

- Update Image for a specific druid CR node in namespace
```
kubectl druid update --cr <cr> --image <image> --namespace <namespace> --node broker
```

- Patch Operation of CR Flags
```
kubectl druid patch --cr <cr> --namespace <namespace> --deleteOrphanPvc true
kubectl druid patch --cr <cr> --namespace <namespace> --rollingDeploy true
```

- Shorthand supported
```
- n for namespace
```

## Note
- Apache®, Apache Druid, Druid® are either registered trademarks or trademarks of the Apache Software Foundation in the United States and/or other countries. This project, kubectl-druid-plugin, is not an Apache Software Foundation project.
