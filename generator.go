package main

import (
	"errors"
	"io/ioutil"

	opav1beta1 "github.com/open-policy-agent/frameworks/constraint/pkg/apis/templates/v1beta1"
	"sigs.k8s.io/yaml"
)


func generate(conf config) (string, error) {
	baseData, err := ioutil.ReadFile(conf.BaseFile)
	if err != nil {
		return "", err
	}
	var tmpl opav1beta1.ConstraintTemplate
	err = yaml.Unmarshal(baseData, &tmpl)
	if err != nil {
		return "", err
	}

	if len(tmpl.Spec.Targets) != len(conf.RegoFiles) {
		return "", errors.New("length mismatch")
	}

	for i, regoFile := range conf.RegoFiles {
		regoData, err := ioutil.ReadFile(regoFile)
		if err != nil {
			return "", err
		}
		tmpl.Spec.Targets[i].Rego = string(regoData)
	}

	buf, err := yaml.Marshal(&tmpl)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
