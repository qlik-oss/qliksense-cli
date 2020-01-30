package qliksense

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func (q *Qliksense) ViewOperatorCrd() {
	io.WriteString(os.Stdout, q.GetCRDString())
}

// this will return crd,deployment,role, rolebinding,serviceaccount for operator
func (q *Qliksense) GetCRDString() string {
	result := ""
	for _, v := range q.getFileList("crd") {
		result = q.getYamlFile(v)
	}
	for _, v := range q.getFileList("crd-deploy") {
		result = result + q.getYamlFile(v)
	}
	return result
}
func (q *Qliksense) getYamlFile(packrFile string) string {
	s, err := q.CrdBox.FindString(packrFile)
	if err != nil {
		fmt.Printf("Cannot read file %s", packrFile)
	}
	return fmt.Sprintln("#soruce: " + packrFile + "\n\n" + s + "\n---")
}
func (q *Qliksense) getFileList(resourceType string) []string {
	var resList []string
	for _, v := range q.CrdBox.List() {
		if strings.Contains(v, resourceType+"/") {
			resList = append(resList, []string{v}...)
		}
	}
	return resList
}
