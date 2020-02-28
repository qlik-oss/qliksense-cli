package main

import (
	"errors"

	qapi "github.com/qlik-oss/sense-installer/pkg/api"
	"github.com/qlik-oss/sense-installer/pkg/qliksense"
	"github.com/spf13/cobra"
)

func pullQliksenseImages(q *qliksense.Qliksense) *cobra.Command {
	opts := &aboutCommandOptions{}

	cmd := &cobra.Command{
		Use:     "pull",
		Short:   "Pull docker images for offline install",
		Example: `qliksense pull`,
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := getAboutCommandGitRef(args)
			if err != nil {
				return err
			}

			if version != "" {
				qConfig := qapi.NewQConfig(q.QliksenseHome)
				if !qConfig.IsRepoExistForCurrent(version) {
					if err := q.FetchQK8s(version); err != nil {
						return err
					}
				} else if err := qConfig.SwitchCurrentCRToVersionAndProfile(version, opts.Profile); err != nil {
					return err
				}
			}

			return q.PullImagesForCurrentCR()
		},
	}
	f := cmd.Flags()
	f.StringVar(&opts.Profile, "profile", "", "Configuration profile")
	return cmd
}

func pushQliksenseImages(q *qliksense.Qliksense) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "push",
		Short:   "Push docker images for offline install",
		Example: `qliksense push`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qConfig := qapi.NewQConfig(q.QliksenseHome)
			if qcr, err := qConfig.GetCurrentCR(); err != nil {
				return err
			} else if registry := qcr.GetImageRegistry(); registry == "" {
				return errors.New("no image registry in config")
			} else {
				return q.PushImagesForCurrentCR()
			}
		},
	}
	return cmd
}
