package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func runGenerator(base, rego, expected string) error {
	tempdir, err := ioutil.TempDir("", "constraint-template-generator-test")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tempdir)
	}()

	basePath := filepath.Join(tempdir, "sample.yaml")
	err = ioutil.WriteFile(basePath, []byte(base), 0644)
	if err != nil {
		return err
	}

	regoPath := filepath.Join(tempdir, "sample.rego")
	err = ioutil.WriteFile(regoPath, []byte(rego), 0644)
	if err != nil {
		return err
	}

	conf := config{
		RegoFiles: []string{regoPath},
		BaseFile:  basePath,
	}

	actual, err := generate(conf)
	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf("expected: %s, but actual: %s", expected, actual)
	}
	return nil
}

func TestGenerator(t *testing.T) {
	err := runGenerator(
		`apiVersion: templates.gatekeeper.sh/v1beta1
kind: ConstraintTemplate
metadata:
  name: k8srequiredlabels
spec:
  targets:
  - target: admission.k8s.gatekeeper.sh
`,
		`package k8srequiredlabels
violation[{"msg": msg, "details": {"missing_labels": missing}}] {
}`,
		`apiVersion: templates.gatekeeper.sh/v1beta1
kind: ConstraintTemplate
metadata:
  creationTimestamp: null
  name: k8srequiredlabels
spec:
  crd:
    spec:
      names: {}
  targets:
  - rego: |-
      package k8srequiredlabels
      violation[{"msg": msg, "details": {"missing_labels": missing}}] {
      }
    target: admission.k8s.gatekeeper.sh
status: {}
`)

	if err != nil {
		t.Error(err)
	}
}
