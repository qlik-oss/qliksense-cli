package main

import (
	"errors"
	"github.com/qlik-oss/sense-installer/pkg/qliksense"
	"github.com/spf13/cobra"
)

func fetchCmd(q *qliksense.Qliksense) *cobra.Command {
	c := &cobra.Command{
		Use:     "fetch",
		Short:   "fetch a release from qliksense-k8s repo",
		Long:    `fetch a release from qliksense-k8s repo`,
		Example: `qliksense fetch <version>`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires a version (i.e. v1.0.0)")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return q.FetchQK8s(args[0])
		},
	}
	return c
}
