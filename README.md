# K8_GoOps
A Kubernetes operator implemented in Golang.

## Installation

### Linux:

> Commands

- `RELEASE_VERSION=<YOUR_VERSION>` Give the latest version (Recommended)  
- `curl -LO "https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk_linux_amd64"`
- `sudo mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk`

## Create the project

```
operator-sdk init --plugins go/v<version> --domain example.com --repo github.com/example/<repo_name> 
```
```
operator-sdk create api --group <group_name> --version v1alpha1 --kind <name_of_CR> --resource --controller
```
