package preflight

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
)

func (qp *QliksensePreflight) CheckDeployment(namespace string, kubeConfigContents []byte) error {
	clientset, _, err := getK8SClientSet(kubeConfigContents, "")
	if err != nil {
		err = fmt.Errorf("Kube config error: %v\n", err)
		return err
	}

	// Deployment check
	qp.P.LogVerboseMessage("Preflight deployment check: \n")
	qp.P.LogVerboseMessage("--------------------------- \n")
	err = qp.checkPfDeployment(clientset, namespace, "deployment-preflight-check")
	if err != nil {
		qp.P.LogVerboseMessage("Preflight Deployment check: FAILED\n")
		return err
	}
	qp.P.LogVerboseMessage("Completed preflight deployment check\n")

	return nil
}

func (qp *QliksensePreflight) CheckService(namespace string, kubeConfigContents []byte) error {
	clientset, _, err := getK8SClientSet(kubeConfigContents, "")
	if err != nil {
		err = fmt.Errorf("unable to create a kubernetes client: %v\n", err)
		return err
	}
	// Service check
	qp.P.LogVerboseMessage("Preflight service check: \n")
	qp.P.LogVerboseMessage("------------------------ \n")
	err = qp.checkPfService(clientset, namespace)
	if err != nil {
		qp.P.LogVerboseMessage("Preflight Service check: FAILED\n")
		return err
	}
	qp.P.LogVerboseMessage("Completed preflight service check\n")
	return nil
}

func (qp *QliksensePreflight) CheckPod(namespace string, kubeConfigContents []byte) error {
	clientset, _, err := getK8SClientSet(kubeConfigContents, "")
	if err != nil {
		err = fmt.Errorf("error: unable to create a kubernetes client: %v\n", err)
		return err
	}
	// Pod check
	qp.P.LogVerboseMessage("Preflight pod check: \n")
	qp.P.LogVerboseMessage("-------------------- \n")
	err = qp.checkPfPod(clientset, namespace)
	if err != nil {
		qp.P.LogVerboseMessage("Preflight Pod check: FAILED\n")
		return err
	}
	qp.P.LogVerboseMessage("Completed preflight pod check\n")
	return nil
}

func (qp *QliksensePreflight) checkPfPod(clientset *kubernetes.Clientset, namespace string) error {
	// create a pod
	podName := "pod-pf-check"
	commandToRun := []string{}

	imageName, err := qp.GetPreflightConfigObj().GetImageName(nginx, true)
	if err != nil {
		return err
	}
	pod, err := qp.createPreflightTestPod(clientset, namespace, podName, imageName, nil, commandToRun)
	if err != nil {
		err = fmt.Errorf("unable to create pod - %v\n", err)
		return err
	}
	defer qp.deletePod(clientset, namespace, podName)

	if err := waitForPod(clientset, namespace, pod); err != nil {
		return err
	}

	qp.P.LogVerboseMessage("Preflight pod creation check: PASSED\n")
	qp.P.LogVerboseMessage("Cleaning up resources...\n")
	return nil
}

func (qp *QliksensePreflight) checkPfService(clientset *kubernetes.Clientset, namespace string) error {
	// creating service
	serviceName := "svc-pf-check"
	pfService, err := qp.createPreflightTestService(clientset, namespace, serviceName)
	if err != nil {
		err = fmt.Errorf("unable to create service - %v\n", err)
		return err
	}
	defer qp.deleteService(clientset, namespace, serviceName)
	_, err = getService(clientset, namespace, pfService.GetName())
	if err != nil {
		err = fmt.Errorf("unable to retrieve service - %v\n", err)
		return err
	}
	qp.P.LogVerboseMessage("Preflight service creation check: PASSED\n")
	qp.P.LogVerboseMessage("Cleaning up resources...\n")
	return nil
}

func (qp *QliksensePreflight) checkPfDeployment(clientset *kubernetes.Clientset, namespace, depName string) error {
	// check if we are able to create a deployment
	imageName, err := qp.GetPreflightConfigObj().GetImageName(nginx, true)
	if err != nil {
		return err
	}
	pfDeployment, err := qp.createPreflightTestDeployment(clientset, namespace, depName, imageName)
	if err != nil {
		err = fmt.Errorf("unable to create deployment - %v\n", err)
		return err
	}
	defer qp.deleteDeployment(clientset, namespace, depName)
	if err := waitForDeployment(clientset, namespace, pfDeployment); err != nil {
		return err
	}
	qp.P.LogVerboseMessage("Preflight Deployment check: PASSED\n")
	qp.P.LogVerboseMessage("Cleaning up resources...\n")
	return nil
}
