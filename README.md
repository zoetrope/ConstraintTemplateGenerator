# ConstraintTemplateGenerator

ConstraintTemplateGenerator is a [Kustomize](https://github.com/kubernetes-sigs/kustomize) plugin 
that generates a ConstraintTemplate by inserting rego files into a base file.

## Install

```console
make install
```

## Usage

Given input files like following:

* kustomization.yaml
    ```yaml
    apiVersion: kustomize.config.k8s.io/v1beta1
    kind: Kustomization
    generators:
        - constraint-template-generator-config.yaml
    ```
* constraint-template-generator-config.yaml
    ```yaml
    kind: ConstraintTemplateGenerator
    apiVersion: "kustomize.cybozu.com/v1"
    metadata:
        name: constraint-template-generator-config
        base: sample.yaml
    regos:
        - sample.rego
    ```
* sample.yaml
    ```yaml
    apiVersion: templates.gatekeeper.sh/v1beta1
    kind: ConstraintTemplate
    metadata:
        name: k8srequiredlabels
    spec:
        crd:
            spec:
            names:
                kind: K8sRequiredLabels
            validation:
                openAPIV3Schema:
                properties:
                    labels:
                    type: array
                    items:
                        type: string
        targets:
            - target: admission.k8s.gatekeeper.sh
    ```
* sample.rego
    ```
    package k8srequiredlabels

    violation[{"msg": msg, "details": {"missing_labels": missing}}] {
        provided := {label | input.review.object.metadata.labels[label]}
        required := {label | label := input.parameters.labels[_]}
        missing := required - provided
        count(missing) > 0
        msg := sprintf("you must provide labels: %v", [missing])
    }
    ```

Specify the directory path of the kustomization.yaml and run the following command.

```console
kustomize build --enable_alpha_plugins ./examples
```

It will generate `ConstraintTemplate` with the required rego added:

```yaml
apiVersion: templates.gatekeeper.sh/v1beta1
kind: ConstraintTemplate
metadata:
  name: k8srequiredlabels
spec:
  crd:
    spec:
      names:
        kind: K8sRequiredLabels
      validation:
        openAPIV3Schema:
          properties:
            labels:
              items:
                type: string
              type: array
  targets:
  - rego: |
      package k8srequiredlabels

      violation[{"msg": msg, "details": {"missing_labels": missing}}] {
        provided := {label | input.review.object.metadata.labels[label]}
        required := {label | label := input.parameters.labels[_]}
        missing := required - provided
        count(missing) > 0
        msg := sprintf("you must provide labels: %v", [missing])
      }
    target: admission.k8s.gatekeeper.sh
```
