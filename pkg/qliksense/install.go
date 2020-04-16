package qliksense

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/qlik-oss/k-apis/pkg/config"
	"sigs.k8s.io/kustomize/api/filesys"

	qapi "github.com/qlik-oss/sense-installer/pkg/api"
)

type InstallCommandOptions struct {
	StorageClass string
	MongoDbUri   string
	RotateKeys   string
}

func (q *Qliksense) InstallQK8s(version string, opts *InstallCommandOptions, keepPatchFiles bool) error {

	// step1: fetch 1.0.0 # pull down qliksense-k8s@1.0.0
	// step2: operator view | kubectl apply -f # operator manifest (CRD)
	// step3: config apply | kubectl apply -f # generates patches (if required) in configuration directory, applies manifest
	// step4: config view | kubectl apply -f # generates Custom Resource manifest (CR)

	// fetch the version
	qConfig := qapi.NewQConfig(q.QliksenseHome)
	if !keepPatchFiles {
		if err := q.DiscardAllUnstagedChangesFromGitRepo(qConfig); err != nil {
			fmt.Printf("error removing temporary changes to the config: %v\n", err)
		}
	}

	qcr, err := qConfig.GetCurrentCR()
	if err != nil {
		fmt.Println("cannot get the current-context cr", err)
		return err
	}

	qcr.SetEULA("yes")
	if opts.MongoDbUri != "" {
		qcr.Spec.AddToSecrets("qliksense", "mongoDbUri", opts.MongoDbUri, "")
	}
	if opts.StorageClass != "" {
		qcr.Spec.StorageClassName = opts.StorageClass
	}
	if opts.RotateKeys != "" {
		qcr.Spec.RotateKeys = opts.RotateKeys
	}
	qConfig.WriteCurrentContextCR(qcr)

	//if the docker pull secret exists on disk, install it in the cluster
	//if it doesn't exist on disk, remove it in the cluster
	if err := installOrRemoveImagePullSecret(qConfig); err != nil {
		return err
	}

	// check if acceptEULA is yes or not
	if !qcr.IsEULA() {
		return errors.New(agreementTempalte + "\n Please do $ qliksense install --acceptEULA=yes\n")
	}

	//CRD will be installed outside of operator
	//install operator controller into the namespace
	fmt.Println("Installing operator controller")
	if operatorControllerString, err := q.getProcessedOperatorControllerString(qcr); err != nil {
		fmt.Println("error extracting/transforming operator controller", err)
		return err
	} else if err := qapi.KubectlApply(operatorControllerString, ""); err != nil {
		fmt.Println("cannot do kubectl apply on operator controller", err)
		return err
	}

	// create patch dependent resoruces
	fmt.Println("Installing resoruces used kuztomize patch")
	if err := q.createK8sResoruceBeforePatch(qcr); err != nil {
		return err
	}

	if qcr.Spec.Git != nil && qcr.Spec.Git.Repository != "" {
		// fetching and applying manifest will be in the operator controller
		// get decrypted cr
		if dcr, err := qConfig.GetDecryptedCr(qcr); err != nil {
			return err
		} else {
			return q.applyCR(dcr)
		}
	}
	if !qcr.IsRepoExist() {
		if err := fetchAndUpdateCR(qConfig, version); err != nil {
			return err
		}
	}

	qcr, err = qConfig.GetCurrentCR()
	if err != nil {
		fmt.Println("cannot get the current-context cr", err)
		return err
	} else if qcr.Spec.GetManifestsRoot() == "" {
		return errors.New("cannot get the manifest root. Use qliksense fetch <version> or qliksense set manifestsRoot")
	}

	// install generated manifests into cluster
	fmt.Println("Installing generated manifests into cluster")

	if dcr, err := qConfig.GetDecryptedCr(qcr); err != nil {
		return err
	} else {
		if IsQliksenseInstalled(dcr.GetName()) {
			return q.UpgradeQK8s(keepPatchFiles)
		}
		if err := q.applyConfigToK8s(dcr); err != nil {
			fmt.Println("cannot do kubectl apply on manifests")
			return err
		} else {
			return q.applyCR(dcr)
		}
	}
}

func (q *Qliksense) getProcessedOperatorControllerString(qcr *qapi.QliksenseCR) (string, error) {
	operatorControllerString := q.GetOperatorControllerString()
	if imageRegistry := qcr.Spec.GetImageRegistry(); imageRegistry != "" {
		return kustomizeForImageRegistry(operatorControllerString, pullSecretName,
			path.Join(qliksenseOperatorImageRepo, qliksenseOperatorImageName),
			path.Join(imageRegistry, qliksenseOperatorImageName))
	}
	return operatorControllerString, nil
}

func installOrRemoveImagePullSecret(qConfig *qapi.QliksenseConfig) error {
	if pullDockerConfigJsonSecret, err := qConfig.GetPullDockerConfigJsonSecret(); err == nil {
		if dockerConfigJsonSecretYaml, err := pullDockerConfigJsonSecret.ToYaml(""); err != nil {
			return err
		} else if err := qapi.KubectlApply(string(dockerConfigJsonSecretYaml), ""); err != nil {
			return err
		}
	} else {
		deleteDockerConfigJsonSecret := qapi.DockerConfigJsonSecret{
			Name: pullSecretName,
		}
		if deleteDockerConfigJsonSecretYaml, err := deleteDockerConfigJsonSecret.ToYaml(""); err != nil {
			return err
		} else if err := qapi.KubectlDelete(string(deleteDockerConfigJsonSecretYaml), ""); err != nil {
			qapi.LogDebugMessage("failed deleting %v, error: %v\n", pullSecretName, err)
		}
	}
	return nil
}

func kustomizeForImageRegistry(resources, dockerConfigJsonSecretName, name, newName string) (string, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	if err := ioutil.WriteFile(filepath.Join(dir, "resources.yaml"), []byte(resources), os.ModePerm); err != nil {
		return "", err
	} else if err := ioutil.WriteFile(filepath.Join(dir, "addImagePullSecrets.yaml"), []byte(fmt.Sprintf(`
apiVersion: builtin
kind: PatchTransformer
metadata:
  name: notImportantHere
patch: '[{"op": "add", "path": "/spec/template/spec/imagePullSecrets", "value": [{"name": "%v"}]}]'
target:
  name: .*-operator
  kind: Deployment
`, dockerConfigJsonSecretName)), os.ModePerm); err != nil {
		return "", err
	} else if err := ioutil.WriteFile(filepath.Join(dir, "kustomization.yaml"), []byte(fmt.Sprintf(`
resources:
- resources.yaml
transformers:
- addImagePullSecrets.yaml
images:
- name: %s
  newName: %s
`, name, newName)), os.ModePerm); err != nil {
		return "", err
	} else if out, err := executeKustomizeBuildForFileSystem(dir, filesys.MakeFsOnDisk()); err != nil {
		return "", err
	} else {
		return string(out), nil
	}
}

func (q *Qliksense) applyCR(cr *qapi.QliksenseCR) error {
	// install operator cr into cluster
	//get the current context cr
	fmt.Println("Install operator CR into cluster")
	r, err := cr.GetString()
	if err != nil {
		return err
	}
	if err := qapi.KubectlApply(r, ""); err != nil {
		fmt.Println("cannot do kubectl apply on operator CR")
		return err
	}
	return nil
}

func (q *Qliksense) createK8sResoruceBeforePatch(qcr *qapi.QliksenseCR) error {
	for svc, nvs := range qcr.Spec.Secrets {
		for _, nv := range nvs {
			if isK8sSecretNeedToCreate(nv) {
				fmt.Println(filepath.Join(qcr.GetK8sSecretsFolder(q.QliksenseHome), svc+".yaml"))
				if secS, err := q.PrepareK8sSecret(filepath.Join(qcr.GetK8sSecretsFolder(q.QliksenseHome), svc+".yaml")); err != nil {
					return err
				} else {
					return qapi.KubectlApply(secS, "")
				}
			}
		}
	}
	return nil
}

func isK8sSecretNeedToCreate(nv config.NameValue) bool {
	return nv.ValueFrom != nil
}
