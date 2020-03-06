package qliksense

import (
	"fmt"

	qapi "github.com/qlik-oss/sense-installer/pkg/api"
)

func (q *Qliksense) UpgradeQK8s(keepPatchFiles bool) error {

	// step1: get CR
	// step2: run kustomize
	// step3: run kubectl apply

	// fetch the version
	qConfig := qapi.NewQConfig(q.QliksenseHome)
	if !keepPatchFiles {
		defer func() {
			if err := q.DiscardAllUnstagedChangesFromGitRepo(qConfig); err != nil {
				fmt.Printf("error removing temporary changes to the config: %v\n", err)
			}
		}()
	}

	qcr, err := qConfig.GetCurrentCR()
	if err != nil {
		fmt.Println("cannot get the current-context cr", err)
		return err
	}
	qcr.Spec.RotateKeys = "no"
	if err := q.applyConfigToK8s(qcr); err != nil {
		fmt.Println("cannot do kubectl apply on manifests")
		return err
	}

	fmt.Println("Install operator CR into cluster")
	r, err := qcr.GetString()
	if err != nil {
		return err
	}
	if err := qapi.KubectlApply(r, ""); err != nil {
		fmt.Println("cannot do kubectl apply on operator CR")
	}
	return nil

}
