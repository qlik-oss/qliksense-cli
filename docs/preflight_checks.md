##Preflight checks
Preflight checks provide pre-installation cluster conformance testing and validation before we install qliksense on the cluster. We gather a suite of conformance tests that can be easily written and run on the target cluster to verify that cluster-specific requirements are met.
The suite consists of a set of `collectors` which run the specifications of every test and `analyzers` which analyze the results of every test run by the collector.
We support the following tests at the moment as part of preflight checks, and the range of the suite will be expanded in future.

Run the following command to view help about the commands supported by preflight at any moment:
```console
$ qliksense preflight
perform preflight checks on the cluster

Usage:
  qliksense preflight [command]

Examples:
qliksense preflight <preflight_check_to_run>

Available Commands:
  all         perform all checks
  dns         perform preflight dns check
  k8s-version check k8s version

Flags:
  -h, --help   help for preflight
```

### DNS check
Run the following command to perform preflight DNS check. We setup a kubernetes deployment and try to reach it as part of establishing DNS connectivity in this check. 
The expected output should be similar to the one shown below.
```console
$ qliksense preflight dns

Preflight DNS check
---------------------
Created deployment "dep-dns-preflight-check"
Created service "svc-dns-pf-check"
Created pod: pf-pod-1
Fetching pod: pf-pod-1
Fetching pod: pf-pod-1
Exec-ing into the container...
Preflight DNS check: PASSED
Completed preflight DNS check
Cleaning up resources...
Deleted pod: pf-pod-1
Deleted service: svc-dns-pf-check
Deleted deployment: dep-dns-preflight-check

```

### Kubernetes version check
We check the version of the target kubernetes cluster and ensure that it falls in the valid range of kubernetes versions that are supported by qliksense. 
The command to run this check and the expected similar output are as shown below:
```console
$ qliksense preflight k8s-version

Preflight kubernetes minimum version check
------------------------------------------
Kubernetes API Server version: v1.15.5
Current K8s Version: 1.15.5
Current 1.15.5 is greater than minimum required version:1.11.0, hence good to go
Preflight minimum kubernetes version check: PASSED
Completed Preflight kubernetes minimum version check

```

### Running all checks
Run the command shown below to execute all preflight checks.
```console
$ qliksense preflight all

Running all preflight checks

Preflight DNS check
-------------------
Created deployment "dep-dns-preflight-check"
Created service "svc-dns-pf-check"
Created pod: pf-pod-1
Fetching pod: pf-pod-1
Fetching pod: pf-pod-1
Exec-ing into the container...
Preflight DNS check: PASSED
Completed preflight DNS check
Cleaning up resources...
Deleted pod: pf-pod-1
Deleted service: svc-dns-pf-check
Deleted deployment: dep-dns-preflight-check

Preflight kubernetes minimum version check
------------------------------------------
Kubernetes API Server version: v1.15.5
Current K8s Version: 1.15.5
Current 1.15.5 is greater than minimum required version:1.11.0, hence good to go
Preflight minimum kubernetes version check: PASSED
Completed Preflight kubernetes minimum version check

All preflight checks have PASSED
Completed running all preflight checks

```