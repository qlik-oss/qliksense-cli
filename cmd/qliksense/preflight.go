package main

import (
	"fmt"
	"log"

	"github.com/qlik-oss/sense-installer/pkg/preflight"

	"github.com/qlik-oss/sense-installer/pkg/qliksense"
	"github.com/spf13/cobra"
)

func preflightCmd(q *qliksense.Qliksense) *cobra.Command {
	var preflightCmd = &cobra.Command{
		Use:     "preflight",
		Short:   "perform preflight checks on the cluster",
		Long:    `perform preflight checks on the cluster`,
		Example: `qliksense preflight <preflight_check_to_run>`,
	}
	return preflightCmd
}

func pfDnsCheckCmd(q *qliksense.Qliksense) *cobra.Command {
	var preflightDnsCmd = &cobra.Command{
		Use:     "dns",
		Short:   "perform preflight dns check",
		Long:    `perform preflight dns check to check DNS connectivity status in the cluster`,
		Example: `qliksense preflight dns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qp := &preflight.QliksensePreflight{Q: q}

			// Preflight DNS check
			fmt.Printf("Preflight DNS check\n")
			fmt.Println("---------------------")
			namespace, kubeConfigContents, err := preflight.InitPreflight()
			if err != nil {
				fmt.Printf("Preflight DNS check FAILED\n")
				log.Fatal(err)
			}
			if err = qp.CheckDns(namespace, kubeConfigContents); err != nil {
				fmt.Println(err)
				fmt.Print("Preflight DNS check FAILED\n")
				log.Fatal()
			}
			return nil
		},
	}
	return preflightDnsCmd
}

func pfK8sVersionCheckCmd(q *qliksense.Qliksense) *cobra.Command {
	var preflightCheckK8sVersionCmd = &cobra.Command{
		Use:     "k8s-version",
		Short:   "check k8s version",
		Long:    `check minimum valid k8s version on the cluster`,
		Example: `qliksense preflight k8s-version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qp := &preflight.QliksensePreflight{Q: q}

			// Preflight Kubernetes minimum version check
			fmt.Printf("Preflight kubernetes minimum version check\n")
			fmt.Println("------------------------------------------")
			namespace, kubeConfigContents, err := preflight.InitPreflight()
			if err != nil {
				fmt.Printf("Preflight kubernetes minimum version check FAILED\n")
				log.Fatal(err)
			}
			if err = qp.CheckK8sVersion(namespace, kubeConfigContents); err != nil {
				fmt.Println(err)
				fmt.Printf("Preflight kubernetes minimum version check FAILED\n")
				log.Fatal()
			}
			return nil
		},
	}
	return preflightCheckK8sVersionCmd
}

func pfAllChecksCmd(q *qliksense.Qliksense) *cobra.Command {
	var preflightAllChecksCmd = &cobra.Command{
		Use:     "all",
		Short:   "perform all checks",
		Long:    `perform all preflight checks on the target cluster`,
		Example: `qliksense preflight all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qp := &preflight.QliksensePreflight{Q: q}

			// Preflight run all checks
			fmt.Printf("Running all preflight checks\n")
			namespace, kubeConfigContents, err := preflight.InitPreflight()
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Running preflight check suite has FAILED...\n")
				log.Fatal()
			}
			qp.RunAllPreflightChecks(namespace, kubeConfigContents)
			return nil

		},
	}
	return preflightAllChecksCmd
}

func pfDeploymentCheckCmd(q *qliksense.Qliksense) *cobra.Command {
	var pfDeploymentCheckCmd = &cobra.Command{
		Use:     "deployment",
		Short:   "perform preflight deploymwnt check",
		Long:    `perform preflight deployment check to ensure that we can create deployments in the cluster`,
		Example: `qliksense preflight deployment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qp := &preflight.QliksensePreflight{Q: q}

			// Preflight deployments check
			fmt.Printf("Preflight deployment check\n")
			fmt.Println("--------------------------")
			namespace, kubeConfigContents, err := preflight.InitPreflight()
			if err != nil {
				fmt.Printf("Preflight deployment check FAILED\n")
				log.Fatal(err)
			}
			if err = qp.CheckDeployment(namespace, kubeConfigContents); err != nil {
				fmt.Println(err)
				fmt.Print("Preflight deploy check FAILED\n")
				log.Fatal()
			}
			return nil
		},
	}
	return pfDeploymentCheckCmd
}

func pfServiceCheckCmd(q *qliksense.Qliksense) *cobra.Command {
	var pfServiceCheckCmd = &cobra.Command{
		Use:     "service",
		Short:   "perform preflight service check",
		Long:    `perform preflight service check to ensure that we are able to create services in the cluster`,
		Example: `qliksense preflight service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qp := &preflight.QliksensePreflight{Q: q}

			// Preflight service check
			fmt.Printf("Preflight service check\n")
			fmt.Println("-----------------------")
			namespace, kubeConfigContents, err := preflight.InitPreflight()
			if err != nil {
				fmt.Printf("Preflight service check FAILED\n")
				log.Fatal(err)
			}
			if err = qp.CheckService(namespace, kubeConfigContents); err != nil {
				fmt.Println(err)
				fmt.Print("Preflight service check FAILED\n")
				log.Fatal()
			}
			return nil
		},
	}
	return pfServiceCheckCmd
}

func pfPodCheckCmd(q *qliksense.Qliksense) *cobra.Command {
	var pfPodCheckCmd = &cobra.Command{
		Use:     "pod",
		Short:   "perform preflight pod check",
		Long:    `perform preflight pod check to ensure we can create pods in the cluster`,
		Example: `qliksense preflight pod`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qp := &preflight.QliksensePreflight{Q: q}

			// Preflight pod check
			fmt.Printf("Preflight pod check\n")
			fmt.Println("--------------------")
			namespace, kubeConfigContents, err := preflight.InitPreflight()
			if err != nil {
				fmt.Printf("Preflight pod check FAILED\n")
				log.Fatal(err)
			}
			if err = qp.CheckPod(namespace, kubeConfigContents); err != nil {
				fmt.Println(err)
				fmt.Print("Preflight pod check FAILED\n")
				log.Fatal()
			}
			return nil
		},
	}
	return pfPodCheckCmd
}

func pfCreateRoleCheckCmd(q *qliksense.Qliksense) *cobra.Command {
	var preflightDnsCmd = &cobra.Command{
		Use:     "create-role",
		Short:   "preflight create role check",
		Long:    `perform preflight role check to ensure we are able to create a role in the cluster`,
		Example: `qliksense preflight create-role`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qp := &preflight.QliksensePreflight{Q: q}

			// Preflight create-role check
			fmt.Printf("Preflight create-role check\n")
			fmt.Println("---------------------------")
			namespace, kubeConfigContents, err := preflight.InitPreflight()
			if err != nil {
				fmt.Printf("Preflight create-role check FAILED\n")
				log.Fatal(err)
			}
			if err = qp.CreateRoleCheck(namespace, kubeConfigContents); err != nil {
				fmt.Println(err)
				fmt.Print("Preflight role-check FAILED\n")
				log.Fatal()
			}
			return nil
		},
	}
	return preflightDnsCmd
}

// preflightCmd.AddCommand(pfMongoCheckCmd(p))
// preflightCmd.AddCommand(pfServiceCheckCmd(p))
// preflightCmd.AddCommand(pfCreateRoleBindingCheckCmd(p))
// preflightCmd.AddCommand(pfCreateServiceAccountCheckCmd(p))
// preflightCmd.AddCommand(pfCreateRBCheckCmd(p))
