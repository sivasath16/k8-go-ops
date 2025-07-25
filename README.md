# K8_GoOps
A Kubernetes operator helps to scale up/down your deployments based on the time(UTC).

## Description
This Go-based operator will scale up/down your deployments based on the time defined in the CRs. You can pass the following properties for in the CR.

```
start: 7 #AM
end: 13 #PM
deployments:
  - name: abc #Eg: Nginx
    namespace: default
replicas: 5
```

| Parameter | Description |
| --- | --- |
| `start` | start time (24hr/format) |
| `end` | end time (24hr/format) |
| `deploments.name` | name of the deployment to be scaled |
| `deployments.namespace` | current namespace of the deployment |
| `replicas` | total no. of replicas within the time period |


## Overview:
K8_GoOps is a Go-based Kubernetes operator built with Kubebuilder that provides time-based automatic scaling for your Kubernetes deployments. Instead of manually scaling deployments or relying on complex cron jobs, this operator allows you to define scaling schedules using Custom Resources (CRs).

## How it Works
The operator continuously monitors your Custom Resources and compares the current UTC time against your defined schedule. When the current time falls within your specified time window, the operator scales your deployments to the desired replica count. Outside of this window, deployments are scaled down to 0 replicas (or a minimum count you specify).

## Features

- `Time-based scaling`: Scale deployments based on UTC time schedules
- `Multiple deployment support`: Scale multiple deployments with a single Custom Resource
- `Namespace-aware`: Support for deployments across different namespaces
- `Flexible configuration`: Easy YAML-based configuration
- `Built with Kubebuilder`: Follows Kubernetes operator best practices
- `Lightweight`: Minimal resource footprint
- `Cloud-native`: Works with any Kubernetes cluster (KIND, Minikube, EKS, GKE, AKS, etc.)

## Pre-requisites

- **Kubernetes cluster**: v1.20+ (KIND, Minikube, or any managed Kubernetes service)
- **Go**: v1.19+ for development
- **Kubebuilder**: v3.0+ for development
- **kubectl**: Latest version
- **Docker**: For building custom images (optional)

## How to use **K8_Ops**

- Clone the Repository
```
git clone `https://github.com/sivasath16/k8-go-ops`
cd k8-go-ops
```

- Containerized Deployment to Cluster
> Build and push the image
```
make docker-build docker-push IMG=<your-registry>/<operator-name>:tag
```
```
make deploy IMG=<your-registry>/<operator-name>:tag
```

> Deploy the Operator to the Cluster
```
make deploy IMG=<your-registry>/<your-operator>:v0.1.0
```

> Install CRDs
```
make install
```

> Apply a Sample Custom Resource
```
kubectl apply -f config/samples/monitor_v1_monitor.yaml
```

> Clean Up (If required)
```
make uninstall
```
```
make undeploy
```

## Contributing

We welcome all developers to contribute to the development, testing, and enhancement of this operator. Whether you're fixing bugs, adding features, or improving documentation, your contributions are greatly appreciated.

**Development Setup**


- Install Kubebuilder

```
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
```
```
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```

> Verify Installation

```
kubebuilder version
```

- Clone the Repository

```
git clone `https://github.com/sivasath16/k8-go-ops`
cd k8-go-ops
```

- Install the CRDs
```
make install
```

- Run the Controller (If you want the controller to keep running, switch to a new terminal for other processes)
```
make run 
```

- Altering API Definitions (If API definitions are changed, this produces new manifests (CR/CRDs))
```
make manifests
```

- Apply CRDs
```
kubectl apply  -f config/crd/bases/monitor.sivasathwik.online_monitors.yaml
```

- Reapply Sample Custom Resource

```
kubectl delete -f config/samples/monitor_v1_monitor.yaml
```
```
kubectl apply -f config/samples/monitor_v1_monitor.yaml
```

## License
This project is currently not licensed. Please contact the repository owner for usage and distribution terms.
