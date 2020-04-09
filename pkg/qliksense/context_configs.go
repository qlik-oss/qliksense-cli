package qliksense

import (
	"fmt"

	"github.com/qlik-oss/k-apis/pkg/config"
	"github.com/robfig/cron/v3"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/tabwriter"

	b64 "encoding/base64"

	ansi "github.com/mattn/go-colorable"
	"github.com/qlik-oss/sense-installer/pkg/api"
	"github.com/ttacon/chalk"
	_ "gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Below are some constants to support qliksense context setup
	QliksenseConfigFile     = "config.yaml"
	QliksenseContextsDir    = "contexts"
	DefaultQliksenseContext = "qlik-default"
	MaxContextNameLength    = 17
	QliksenseSecretsDir     = "secrets"

	imageRegistryConfigKey     = "imageRegistry"
	pullSecretName             = "artifactory-docker-secret"
	qliksenseOperatorImageRepo = "qlik-docker-oss.bintray.io"
	qliksenseOperatorImageName = "qliksense-operator"
)

// SetSecrets - set-secrets <key>=<value> commands
func (q *Qliksense) SetSecrets(args []string, isSecretSet bool) error {
	qConfig := api.NewQConfig(q.QliksenseHome)
	qliksenseCR, err := qConfig.GetCurrentCR()
	if err != nil {
		return err
	}

	// Metadata name in qliksense CR is the name of the current context
	api.LogDebugMessage("Current context: %s", qliksenseCR.GetName())
	encryptionKey, err := qConfig.GetEncryptionKeyForCurrent()
	if err != nil {
		return err
	}
	resultArgs, err := api.ProcessConfigArgs(args)
	if err != nil {
		return err
	}
	for _, ra := range resultArgs {
		api.LogDebugMessage("value args to be encrypted: %s", ra.Value)
		if err := q.processSecret(ra, encryptionKey, qliksenseCR, isSecretSet); err != nil {
			return err
		}
	}
	// write modified content into context-yaml
	return qConfig.WriteCR(qliksenseCR)
}

func (q *Qliksense) processSecret(ra *api.ServiceKeyValue, encryptionKey string, qliksenseCR *api.QliksenseCR, isSecretSet bool) error {
	cipherText, e2 := api.EncryptData([]byte(ra.Value), encryptionKey)
	if e2 != nil {
		return e2
	}
	secretName := ""
	if isSecretSet {
		secretFolder := qliksenseCR.GetK8sSecretsFolder(q.QliksenseHome)
		secretFileName := filepath.Join(secretFolder, ra.SvcName+".yaml")

		secretName = fmt.Sprintf("%s-%s-%s", qliksenseCR.GetName(), ra.SvcName, "senseinstaller")
		api.LogDebugMessage("Constructed secret name: %s", secretName)

		k8sSecret := v1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: secretName,
			},
			Type: v1.SecretTypeOpaque,
		}

		if !api.DirExists(secretFolder) {
			if err := os.MkdirAll(secretFolder, os.ModePerm); err != nil {
				err = fmt.Errorf("Not able to create %s dir: %v", secretFolder, err)
				log.Println(err)
				return err
			}
		}
		// if read from file errors out, we can ignore it here
		_ = api.ReadFromFile(&k8sSecret, secretFileName)
		if k8sSecret.Data == nil {
			k8sSecret.Data = map[string][]byte{}
		}
		// v1.Secret does enconding, so no need to encode again
		k8sSecret.Data[ra.Key] = []byte(cipherText)

		// Write secret to file
		k8sSecretBytes, err := api.K8sSecretToYaml(k8sSecret)
		if err != nil {
			api.LogDebugMessage("Error while converting K8s secret to yaml")
			return err
		}
		if err = ioutil.WriteFile(secretFileName, k8sSecretBytes, os.ModePerm); err != nil {
			api.LogDebugMessage("Error while writing K8s secret to file")
			return err
		}
		api.LogDebugMessage("Created a Kubernetes secret")
	}
	base64EncodedSecret := b64.StdEncoding.EncodeToString([]byte(cipherText))
	// write into CR the keyref of the secret
	qliksenseCR.Spec.AddToSecrets(ra.SvcName, ra.Key, base64EncodedSecret, secretName)
	return nil
}

// SetConfigs - set-configs <key>=<value> commands
func (q *Qliksense) SetConfigs(args []string) error {
	// retieve current context from config.yaml
	qConfig := api.NewQConfig(q.QliksenseHome)
	qliksenseCR, err := qConfig.GetCurrentCR()
	if err != nil {
		return err
	}

	resultArgs, err := api.ProcessConfigArgs(args)
	if err != nil {
		return err
	}
	for _, ra := range resultArgs {
		qliksenseCR.Spec.AddToConfigs(ra.SvcName, ra.Key, ra.Value)
	}
	// write modified content into context.yaml
	return qConfig.WriteCR(qliksenseCR)
}

func caseInsenstiveFieldByName(v reflect.Value, name string) reflect.Value {
	name = strings.ToLower(name)
	return v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == name })
}

func validateCR(key string, keySub string, value string, crSpec *api.QliksenseCR) (bool, *api.QliksenseCR) {
	cr := reflect.ValueOf(crSpec.Spec)
	keyValid := caseInsenstiveFieldByName(reflect.Indirect(cr), key)
	if !keyValid.IsValid() {
		//not in main spec
		fmt.Println(key, "is an invalid key")
		return false, crSpec
	} else if keySub == "" {
		if key == "rotatekeys" {
			if _, err := validateInput(value); err != nil {
				return false, crSpec
			}
		}
	}
	// checks if it is git or gitops
	if keySub != "" {
		if !keyValid.IsNil() {
			if !caseInsenstiveFieldByName(reflect.Indirect(keyValid), keySub).IsValid() {
				fmt.Println(keySub, "is an invalid key")
				return false, crSpec
			} else {
				// verify gitops enabled and gitops schedule
				switch keySub {
				case "schedule":
					if _, err := cron.ParseStandard(value); err != nil {
						fmt.Println("Please enter string with standard cron scheduling syntax ")
						return false, crSpec
					}
				case "enabled":
					if !strings.EqualFold(value, "yes") && !strings.EqualFold(value, "no") {
						fmt.Println("Please use yes or no for key enabled")
						return false, crSpec
					}
				}
			}
		} else {
			switch key {
			case "gitops":
				crSpec.Spec.GitOps = &config.GitOps{}
			case "git":
				crSpec.Spec.Git = &config.Repo{}
			}
		}
	}
	return true, crSpec
}

// SetOtherConfigs - set profile/storageclassname/git.repository/manifestRoot commands
func (q *Qliksense) SetOtherConfigs(args []string) error {
	// retieve current context from config.yaml
	qConfig := api.NewQConfig(q.QliksenseHome)
	qliksenseCR, err := qConfig.GetCurrentCR()
	if err != nil {
		return err
	}

	// modify appropriate fields
	if len(args) == 0 {
		err := fmt.Errorf("No args were provided. Please provide args to configure the current context")
		log.Println(err)
		return err
	}

	for _, arg := range args {
		argsString := strings.Split(arg, "=")
		key := strings.ToLower(argsString[0])
		value := argsString[1]
		// check if key is for git or gitops (sub objects)
		keySplit := strings.Split(key, ".")
		key = keySplit[0]
		keySub := ""

		if len(keySplit) == 2 {
			keySub = strings.ToLower(keySplit[1])
		}

		valid := true
		valid, qliksenseCR = validateCR(key, keySub, value, qliksenseCR)
		field := caseInsenstiveFieldByName(reflect.Indirect(reflect.ValueOf(qliksenseCR.Spec)), key)
		if !valid {
			err := fmt.Errorf("Please enter one of: profile, storageClassName,rotateKeys, manifestRoot, git.repository or gitops arguments to configure the current context")
			return err
		} else if strings.EqualFold("", keySub) {
			// set spec for everything excluding git and gitops
			if field.CanSet() {
				field.SetString(value)
			}
		} else {
			// set spec for git or gitops
			subField := caseInsenstiveFieldByName(reflect.Indirect(field), keySub)
			if subField.CanSet() {
				subField.SetString(value)
			}
		}

		fmt.Println(chalk.Green.Color("Successfully added to Custom Resource Spec"))
	}
	// write modified content into context.yaml
	return qConfig.WriteCR(qliksenseCR)
}

// SetContextConfig - set the context for qliksense kubernetes resources to live in
func (q *Qliksense) SetContextConfig(args []string) error {
	if len(args) == 1 {
		err := q.SetUpQliksenseContext(args[0])
		if err != nil {
			return err
		}
	} else {
		err := fmt.Errorf("Please provide a name to configure the context with")
		log.Println(err)
		return err
	}
	return nil
}

func (q *Qliksense) ListContextConfigs() error {
	qliksenseConfigFile := filepath.Join(q.QliksenseHome, QliksenseConfigFile)
	var qliksenseConfig api.QliksenseConfig
	if err := api.ReadFromFile(&qliksenseConfig, qliksenseConfigFile); err != nil {
		log.Println(err)
		return err
	}
	out := ansi.NewColorableStdout()
	w := tabwriter.NewWriter(out, 5, 8, 0, '\t', 0)
	fmt.Fprintln(w, chalk.Underline.TextStyle("Context Name"), "\t", chalk.Underline.TextStyle("CR File Location"))
	w.Flush()
	if len(qliksenseConfig.Spec.Contexts) > 0 {
		for _, cont := range qliksenseConfig.Spec.Contexts {
			fmt.Fprintln(w, cont.Name, "\t", qliksenseConfig.GetCRFilePath(cont.Name), "\t")
		}
		w.Flush()
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, chalk.Bold.TextStyle("Current Context : "), qliksenseConfig.Spec.CurrentContext)
	} else {
		fmt.Fprintln(out, "No Contexts Available")
	}
	return nil
}

func (q *Qliksense) DeleteContextConfig(args []string) error {
	if len(args) == 1 {
		qliksenseConfigFile := filepath.Join(q.QliksenseHome, QliksenseConfigFile)
		var qliksenseConfig api.QliksenseConfig
		api.ReadFromFile(&qliksenseConfig, qliksenseConfigFile)
		out := ansi.NewColorableStdout()
		switch args[0] {
		case qliksenseConfig.Spec.CurrentContext:
			fmt.Fprintln(out, chalk.Yellow.Color("Please switch contexts to be able to delete this context."))
			err := fmt.Errorf(chalk.Red.Color("Cannot delete current context - %s"), chalk.White.Color(chalk.Bold.TextStyle(qliksenseConfig.Spec.CurrentContext)))
			return err
		case DefaultQliksenseContext:
			err := fmt.Errorf(chalk.Red.Color("Cannot delete default qliksense context"))
			return err
		default:
			qliksenseContextsDir1 := filepath.Join(q.QliksenseHome, QliksenseContextsDir)
			qliksenseContextFile := filepath.Join(qliksenseContextsDir1, args[0])
			qliksenseSecretsDir1 := filepath.Join(q.QliksenseHome, QliksenseSecretsDir, QliksenseContextsDir)
			qliksenseSecretsFile := filepath.Join(qliksenseSecretsDir1, args[0])
			if err := os.RemoveAll(qliksenseContextFile); err != nil {
				err = fmt.Errorf("Not able to delete %s dir: %v", qliksenseContextsDir1, err)
				log.Println(err)
				return err
			} else if err := os.RemoveAll(qliksenseSecretsFile); err != nil {
				err = fmt.Errorf("No Secrets Folder Detected")
				log.Println(err)
				return err
			} else {
				currentLength := len(qliksenseConfig.Spec.Contexts)
				if currentLength > 0 {
					temp := qliksenseConfig.Spec.Contexts
					qliksenseConfig.Spec.Contexts = nil
					for _, ctx := range temp {
						if ctx.Name != args[0] {
							qliksenseConfig.Spec.Contexts = append(qliksenseConfig.Spec.Contexts, api.Context{
								Name:   ctx.Name,
								CrFile: ctx.CrFile,
							})
						}
					}
					newLength := len(qliksenseConfig.Spec.Contexts)
					if currentLength != newLength {
						api.WriteToFile(&qliksenseConfig, qliksenseConfigFile)
						fmt.Fprintln(out, chalk.Yellow.Color(chalk.Underline.TextStyle("Warning: Active resources may still be running in-cluster")))
						fmt.Fprintln(out, chalk.Green.Color("Successfully deleted context: "), chalk.Bold.TextStyle(args[0]))
					} else {
						err := fmt.Errorf(chalk.Red.Color("Context not found"))
						return err
					}
				}
			}
		}
	} else {
		err := fmt.Errorf("Please provide a context as an argument to delete")
		log.Println(err)
		return err
	}
	return nil
}

// SetUpQliksenseDefaultContext - to setup dir structure for default qliksense context
func (q *Qliksense) SetUpQliksenseDefaultContext() error {
	if api.FileExists(filepath.Join(q.QliksenseHome, "config.yaml")) {
		qliksenseConfig := api.NewQConfig(q.QliksenseHome)
		if qliksenseConfig.IsContextExist(DefaultQliksenseContext) {
			return nil
		}
	}
	return q.SetUpQliksenseContext(DefaultQliksenseContext)
}

// SetUpQliksenseContext - to setup qliksense context
func (q *Qliksense) SetUpQliksenseContext(contextName string) error {
	if contextName == "" {
		err := fmt.Errorf("Please enter a non-empty context-name")
		log.Println(err)
		return err
	}
	// check the length of the context name entered by the user, it should not exceed 17 chars
	if len(contextName) > MaxContextNameLength {
		err := fmt.Errorf("Please enter a context-name with utmost 17 characters")
		log.Println(err)
		return err
	}

	qliksenseConfigFile := filepath.Join(q.QliksenseHome, QliksenseConfigFile)
	qliksenseConfig := api.NewQConfigEmpty(q.QliksenseHome)

	if !api.FileExists(qliksenseConfigFile) {
		qliksenseConfig.AddBaseQliksenseConfigs(contextName)
	} else {
		if err := api.ReadFromFile(qliksenseConfig, qliksenseConfigFile); err != nil {
			log.Println(err)
			return err
		}
	}

	if qliksenseConfig.IsContextExist(contextName) {
		qliksenseConfig.Spec.CurrentContext = contextName
		return qliksenseConfig.Write()
	}
	qliksenseCR := &api.QliksenseCR{}
	qliksenseCR.AddCommonConfig(contextName)
	qliksenseConfig.Spec.CurrentContext = contextName
	if err := qliksenseConfig.CreateOrWriteCrAndContext(qliksenseCR); err != nil {
		return err
	}

	// set the encrypted default mongo
	return q.SetSecrets([]string{`qliksense.mongoDbUri="mongodb://qlik-default-mongodb:27017/qliksense?ssl=false"`}, false)
}

func validateInput(input string) (string, error) {
	var err error
	validInputs := []string{"yes", "no", "None"}
	isValid := false
	for _, elem := range validInputs {
		if input == elem {
			isValid = true
			break
		}
	}
	if !isValid {
		err = fmt.Errorf("Please enter one of: yes, no or None")
		log.Println(err)

	}
	return input, err
}

// PrepareK8sSecret targetFile contains base64 encoded value of encrypted value.
// this method decodes and decrypts the secret value in the secret.yaml file and returns a B64encoded string
func (q *Qliksense) PrepareK8sSecret(targetFile string) (string, error) {
	// check if targetFile exists
	if !api.FileExists(targetFile) {
		err := fmt.Errorf("Target file does not exist in the path provided")
		log.Println(err)
		return "", err
	}
	qConfig := api.NewQConfig(q.QliksenseHome)
	encryptionKey, err := qConfig.GetEncryptionKeyForCurrent()
	if err != nil {
		return "", err
	}

	// read the target file
	k8sSecret, err := readTargetfile(targetFile)
	if err != nil {
		return "", err
	}
	// retrieve value from data section
	k8sSecret1, err := api.K8sSecretFromYaml(k8sSecret)
	if err != nil {
		return "", err
	}
	dataMap := k8sSecret1.Data
	var resultMap = make(map[string][]byte)
	for k, v := range dataMap {
		//k8s secrets has already base64 decoed value
		decryptedString, err := api.DecryptData(v, encryptionKey)
		if err != nil {
			err := fmt.Errorf("Not able to decrypt message: %v", err)
			return "", err
		}
		resultMap[k] = []byte(decryptedString)
	}

	// putting the above map back into the k8sSecret struct
	k8sSecret1.Data = resultMap
	k8sSecretBytes, err := api.K8sSecretToYaml(k8sSecret1)
	if err != nil {
		return "", err
	}
	return string(k8sSecretBytes), nil
}

func readTargetfile(targetFile string) ([]byte, error) {
	k8sSecret, err := ioutil.ReadFile(targetFile)
	if err != nil {
		err := fmt.Errorf("Unable to read the targetFile")
		log.Println(err)
		return nil, err
	}
	return k8sSecret, nil
}

func (q *Qliksense) SetImageRegistry(registry, pushUsername, pushPassword, pullUsername, pullPassword string) error {
	qConfig := api.NewQConfig(q.QliksenseHome)
	qliksenseCR, err := qConfig.GetCurrentCR()
	if err != nil {
		return err
	}
	if pushUsername != "" {
		if err := qConfig.SetPushDockerConfigJsonSecret(&api.DockerConfigJsonSecret{
			Uri:      registry,
			Username: pushUsername,
			Password: pushPassword,
		}); err != nil {
			return err
		} else if err := qConfig.SetPullDockerConfigJsonSecret(&api.DockerConfigJsonSecret{
			Name:     pullSecretName,
			Uri:      registry,
			Username: pullUsername,
			Password: pullPassword,
			Email:    pullUsername,
		}); err != nil {
			return err
		}
	} else if err := qConfig.DeletePushDockerConfigJsonSecret(); err != nil && !os.IsNotExist(err) {
		return err
	} else if err := qConfig.DeletePullDockerConfigJsonSecret(); err != nil && !os.IsNotExist(err) {
		return err
	}

	qliksenseCR.Spec.AddToConfigs("qliksense", imageRegistryConfigKey, registry)
	return qConfig.WriteCR(qliksenseCR)
}

func (q *Qliksense) SetEulaAccepted() error {
	qConfig := api.NewQConfig(q.QliksenseHome)
	qcr, err := qConfig.GetCurrentCR()
	if err != nil {
		return err
	}
	if !qcr.IsEULA() {
		qcr.SetEULA("yes")
		return qConfig.WriteCurrentContextCR(qcr)
	}
	return nil
}
