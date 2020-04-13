package preflight

import (
	"fmt"
	"strings"

	qapi "github.com/qlik-oss/sense-installer/pkg/api"
)

const (
	mongo = "mongo"
)

func (qp *QliksensePreflight) CheckMongo(kubeConfigContents []byte, namespace, mongodbUrl string) error {
	fmt.Printf("Preflight mongodb check: \n")

	if mongodbUrl == "" {
		// infer mongoDbUrl from currentCR
		fmt.Println("MongoDbUri is empty, infer from CR")
		qConfig := qapi.NewQConfig(qp.Q.QliksenseHome)
		var currentCR *qapi.QliksenseCR

		var err error
		qConfig.SetNamespace(namespace)
		currentCR, err = qConfig.GetCurrentCR()
		if err != nil {
			fmt.Printf("Unable to retrieve current CR: %v\n", err)
			return err
		}
		decryptedCR, err := qConfig.GetDecryptedCr(currentCR)
		mongodbUrl = decryptedCR.Spec.GetFromSecrets("qliksense", "mongoDbUri")
	}

	fmt.Printf("mongodbUrl: %s\n", mongodbUrl)
	if err := qp.mongoConnCheck(kubeConfigContents, namespace, mongodbUrl); err != nil {
		return err
	}
	fmt.Println("Completed preflight mongodb check")
	return nil
}

func (qp *QliksensePreflight) mongoConnCheck(kubeConfigContents []byte, namespace, mongodbUrl string) error {
	clientset, _, err := getK8SClientSet(kubeConfigContents, "")
	if err != nil {
		err = fmt.Errorf("error: unable to create a kubernetes client: %v\n", err)
		fmt.Println(err)
		return err
	}
	// create a pod
	podName := "pf-mongo-pod"
	commandToRun := []string{"sh", "-c", "sleep 10;mongo " + mongodbUrl}
	imageName, err := qp.GetPreflightConfigObj().GetImageName(mongo, true)
	if err != nil {
		return err
	}
	mongoPod, err := createPreflightTestPod(clientset, namespace, podName, imageName, commandToRun)
	if err != nil {
		err = fmt.Errorf("error: unable to create pod : %v\n", err)
		return err
	}
	defer deletePod(clientset, namespace, podName)

	if err := waitForPod(clientset, namespace, mongoPod); err != nil {
		return err
	}
	if len(mongoPod.Spec.Containers) == 0 {
		err := fmt.Errorf("error: there are no containers in the pod- %v\n", err)
		fmt.Println(err)
		return err
	}
	waitForPodToDie(clientset, namespace, mongoPod)
	logStr, err := getPodLogs(clientset, mongoPod)
	if err != nil {
		err = fmt.Errorf("error: unable to execute mongo check in the cluster: %v\n", err)
		fmt.Println(err)
		return err
	}

	stringToCheck := "Implicit session:"
	if strings.Contains(logStr, stringToCheck) {
		fmt.Println("Preflight mongo check: PASSED")
	} else {
		err = fmt.Errorf("Expected response not found\n")
		return err
	}
	return nil
}
