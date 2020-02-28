package main

import (
	"github.com/qlik-oss/sense-installer/pkg/qliksense"
	"github.com/spf13/cobra"
)

func configCmd(q *qliksense.Qliksense) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "do operations on/around CR",
		Long:  `do operations on/around CR`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return q.ConfigViewCR()
		},
	}
	return configCmd
}

func configApplyCmd(q *qliksense.Qliksense) *cobra.Command {
	c := &cobra.Command{
		Use:     "apply",
		Short:   "generate the patchs and apply manifests to k8s",
		Long:    `generate patches based on CR and apply manifests to k8s`,
		Example: `qliksense config apply`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return q.ConfigApplyQK8s()
		},
	}
	return c
}

func configViewCmd(q *qliksense.Qliksense) *cobra.Command {
	c := &cobra.Command{
		Use:     "view",
		Short:   "view the qliksense operator CR",
		Long:    `display the operator CR, that has been created for the current context`,
		Example: `qliksense config view`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return q.ConfigViewCR()
		},
	}
	return c
}
