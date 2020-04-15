# (WIP) Qlik Sense on Kubernetes installation and operations CLI

## Documentation

To learn more about Qlik Sense on Kubernetes CLI go to https://qlik-oss.github.io/sense-installer/

## About

The QSEoK CLI (qliksense) provides an imperative interface to many of the configurations that need to be applied against the declarative structure described in [qliksense-k8s](https://github.com/qlik-oss/qliksense-k8s).

This is a technology preview that uses Qlik modified [kustomize](https://github.com/qlik-oss/kustomize) to kubernetes manifests of the versions of the [qliksense-k8s](https://github.com/qlik-oss/qliksense-k8s) repository.

For each version of a qliksense edge build there should be a corresponding release in [qliksense-k8s] repository under [releases](https://github.com/qlik-oss/qliksense-k8s/releases)

### Future Direction

- More operations:
  - Expand preflight checks
  - backup/restore operations
  - fully support airgap installation of QSEoK
  - restore unwanted deletion of kubernetes resources
