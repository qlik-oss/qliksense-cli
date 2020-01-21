package main

import (
	"fmt"
	"github.com/qlik-oss/sense-installer/pkg/qliksense"
	"github.com/spf13/cobra"
	"strings"
)

func porter(q *qliksense.Qliksense) *cobra.Command {
	return &cobra.Command{
		Use:                "porter",
		Short:              "Execute a porter command",
		DisableFlagParsing: true,
		RunE: func(cobCmd *cobra.Command, args []string) error {
			var (
				err error
			)
			if _, err = q.CallPorter(args,
				func(x string) (out *string) {
					out = new(string)
					*out = strings.ReplaceAll(x, "porter", "qliksense porter")
					fmt.Println(*out)
					return
				}); err != nil {
				return err
			}
			return nil
		},
	}
}
