package main

import (
	"fmt"

	"github.com/qlik-oss/sense-installer/pkg/qliksense"
	"github.com/spf13/cobra"
)

func installCmd(q *qliksense.Qliksense) *cobra.Command {
	opts := &qliksense.InstallCommandOptions{}
	keepPatchFiles, pull, push := false, false, false
	c := &cobra.Command{
		Use:     "install",
		Short:   "install a qliksense release",
		Long:    `install a qliksense release`,
		Example: `qliksense install <version> #if no version provides, expect manifestsRoot is set somewhere in the file system`,
		RunE: func(cmd *cobra.Command, args []string) error {
			version := ""
			if len(args) != 0 {
				version = args[0]
			}
			if err := validatePullPushFlagsOnInstall(q, pull, push); err != nil {
				return err
			}
			if pull {
				fmt.Println("Pulling images...")
				if err := q.PullImages(version, ""); err != nil {
					return err
				}
			}
			if push {
				fmt.Println("Pushing images...")
				if err := q.PushImagesForCurrentCR(); err != nil {
					return err
				}
			}
			return q.InstallQK8s(version, opts, keepPatchFiles)
		},
	}

	f := c.Flags()
	f.StringVarP(&opts.StorageClass, "storageClass", "s", "", "Storage class for qliksense")
	f.StringVarP(&opts.MongoDbUri, "mongoDbUri", "m", "", "mongoDbUri for qliksense (i.e. mongodb://qlik-default-mongodb:27017/qliksense?ssl=false)")
	f.StringVarP(&opts.RotateKeys, "rotateKeys", "r", "", "Rotate JWT keys for qliksense (yes:rotate keys/ no:use exising keys from cluster/ None: use default EJSON_KEY from env")
	f.BoolVar(&keepPatchFiles, keepPatchFilesFlagName, keepPatchFiles, keepPatchFilesFlagUsage)
	f.BoolVarP(&pull, pullFlagName, pullFlagShorthand, pull, pullFlagUsage)
	f.BoolVarP(&push, pushFlagName, pushFlagShorthand, push, pushFlagUsage)
	return c
}

func validatePullPushFlagsOnInstall(q *qliksense.Qliksense, pull, push bool) error {
	if pull && !push {
		fmt.Printf("WARNING: pulling images without pushing them")
	}
	if push {
		if err := ensureImageRegistrySetInCR(q); err != nil {
			return err
		}
	}
	return nil
}
