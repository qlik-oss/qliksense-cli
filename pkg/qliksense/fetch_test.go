package qliksense

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	qapi "github.com/qlik-oss/sense-installer/pkg/api"
)

func TestFetchAndUpdateCR(t *testing.T) {
	tempHome, _ := ioutil.TempDir("", "")

	q := &Qliksense{
		QliksenseHome: tempHome,
	}
	q.SetUpQliksenseContext("test1")
	qConfig := qapi.NewQConfig(tempHome)
	if err := fetchAndUpdateCR(qConfig, "v0.0.2"); err != nil {
		t.Log(err)
		t.FailNow()
	}

	actualCrFile := filepath.Join(tempHome, "contexts", "test1", "test1.yaml")
	cr := &qapi.QliksenseCR{}
	if err := qapi.ReadFromFile(cr, actualCrFile); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if cr.Spec.ManifestsRoot != "contexts/test1/qlik-k8s/v0.0.2" {
		t.Log("actual path: " + cr.Spec.ManifestsRoot + ", expected path: contexts/test1/qlik-k8s/v0.0.2")
		t.FailNow()
	}
}
