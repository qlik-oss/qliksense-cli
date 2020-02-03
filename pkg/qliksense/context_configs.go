package qliksense

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qlik-oss/k-apis/config"
	"github.com/qlik-oss/sense-installer/pkg/api"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	QliksenseConfigHome        = "/.qliksense"
	QliksenseConfigContextHome = "/.qliksense/contexts"

	QliksenseConfigApiVersion = "config.qlik.com/v1"
	QliksenseConfigKind       = "QliksenseConfig"
	QliksenseMetadataName     = "QliksenseConfigMetadata"

	QliksenseContextApiVersion    = "qlik.com/v1"
	QliksenseContextKind          = "Qliksense"
	QliksenseContextLabel         = "v1.0.0"
	QliksenseContextManifestsRoot = "/Usr/ddd/my-k8-repo/manifests"
)

// ReadQliksenseContextConfig is exported
func ReadQliksenseContextConfig(qliksenseCR *api.QliksenseCR, fileName string) {
	log.Debugf("Reading file %s", fileName)
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading from source: %s\n", err)
	}
	if err = yaml.Unmarshal([]byte(yamlFile), qliksenseCR); err != nil {
		log.Fatalf("Error when parsing from source: %s\n", err)
	}
}

// WriteQliksenseContextConfigToFile is exported
func WriteQliksenseContextConfigToFile(qliksenseConfig *api.QliksenseConfig, qliksenseCR *api.QliksenseCR, qliksenseFile string) {
	log.Debug("Entry: WriteQliksenseContextConfigToFile()")
	if qliksenseCR != nil {
		log.Debug("This action is about writing to a context file")
		if !FileExists(qliksenseFile) {
			log.Debugf("File %s doesnt exist, creating it now...", qliksenseFile)
			file, err := os.OpenFile(qliksenseFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				log.Debug("There was an error creating the file: %s, %v", qliksenseFile, err)
				panic(err)
			}
			log.Debugf("File created: %s", qliksenseFile)
			defer file.Close()

			log.Debugf("Adding CommonConfig to %s", qliksenseFile)
			// infer context name from the filename path
			contextName := strings.Replace(filepath.Base(qliksenseFile), ".yaml", "", 1)

			qliksenseCR1 := AddCommonConfig(*qliksenseCR, contextName)
			log.Debug("Added CommonConfig to %s", qliksenseFile)
			x, err := yaml.Marshal(qliksenseCR1)
			if err != nil {
				log.Fatalf("An error occurred during marshalling CR: %v", err)
			}
			log.Debugf("Marshalled yaml:\n%s\nWriting to file...", x)

			numBytes, err := file.Write(x)
			if err != nil {
				panic(err)
			}
			log.Debugf("wrote %d bytes\n", numBytes)
			log.Debugf("Wrote Struct into %s", qliksenseFile)
		} else {
			log.Debug("File %s already exists ", qliksenseFile)
		}
	} else {
		log.Debug("This section is about writing into the base config file %s", qliksenseFile)
	}
}

// AddCommonConfig is exported
func AddCommonConfig(qliksenseCR api.QliksenseCR, contextName string) api.QliksenseCR {
	log.Debug("Entry: addCommonConfig()")
	qliksenseCR.ApiVersion = QliksenseContextApiVersion
	qliksenseCR.Kind = QliksenseContextKind
	if qliksenseCR.Metadata.Name == "" {
		qliksenseCR.Metadata.Name = contextName
	}
	qliksenseCR.Metadata.Labels = map[string]string{}
	qliksenseCR.Metadata.Labels["Version"] = QliksenseContextLabel
	qliksenseCR.Spec = &config.CRSpec{}
	qliksenseCR.Spec.ManifestsRoot = QliksenseContextManifestsRoot
	log.Debug("Exit: addCommonConfig()")
	return qliksenseCR
}

// AddBaseQliksenseConfigs is exported
func AddBaseQliksenseConfigs(qliksenseConfig api.QliksenseConfig, defaultQliksenseContext string) api.QliksenseConfig {
	log.Debug("Entry: AddBaseQliksenseConfigs()")
	qliksenseConfig.ApiVersion = QliksenseConfigApiVersion
	qliksenseConfig.Kind = QliksenseConfigKind
	qliksenseConfig.Metadata.Name = QliksenseMetadataName
	if defaultQliksenseContext != "" {
		qliksenseConfig.Spec.CurrentContext = defaultQliksenseContext
	}
	log.Debug("Exit: AddBaseQliksenseConfigs()")
	return qliksenseConfig
}

func setOtherConfigs(q *Qliksense) error {
	return nil
}

// FileExists is exported
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Debug("File does not exist")
		return false
	}
	log.Debug("Either File exists OR a different error occurred")
	return !info.IsDir()
}
