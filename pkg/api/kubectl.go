package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func KubectlApply(manifests string) error {
	return kubectlOperation(manifests, "apply")
}

func KubectlDelete(manifests string) error {
	return kubectlOperation(manifests, "delete")
}

func kubectlOperation(manifests string, oprName string) error {
	tempYaml, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Println("cannot create file ", err)
		return err
	}
	tempYaml.WriteString(manifests)

	var cmd *exec.Cmd
	if oprName == "apply" {
		cmd = exec.Command("kubectl", oprName, "-f", tempYaml.Name(), "--validate=false")
	} else {
		cmd = exec.Command("kubectl", oprName, "-f", tempYaml.Name())
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("kubectl apply failed with %s\n", err)
		fmt.Println("temp CRD file: " + tempYaml.Name())
		return err
	}
	os.Remove(tempYaml.Name())
	return nil
}
