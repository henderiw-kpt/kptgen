# kptgen

kptgen is a tool to help build a kpt package. 

kptgen uses a set of operations to build the kpt package, such as:
- kptgen copy: allows to copy files from a certain directory (used only for CRDs for now)
- kptgen init: initializes the kpt package with a set of KRM fn-config templates, that can be used in the apply operations
- kptgen apply pod: add the pod resources to the kpt package (deployment, statefullset, serviceaccount, services, roles, rolebindings, ...)
- kptgen apply clusterrole: add a clusterrole to the kpt package
- kptgen apply namespace: add a namespace to the kpt package
- kptgen apply webhook: adds the webhook, cert, service resources to the kpt packge and mutates the deployment/statefulset with the volumes/volumemounts for the certificate

## usage example

initialize the kpt package

```
export KPT_BLUEPRINT_DIR=./blueprint/admin3
mkdir -p ${KPT_BLUEPRINT_DIR}
kpt init ${KPT_BLUEPRINT_DIR}
set the namespace in the KptFile -> to be automated
```

```
kpt pkg clone GIT_REPO_URL ${KPT_BLUEPRINT_DIR} // clone the fn-config templates
kptgen copy SOURCE_DIR ${KPT_BLUEPRINT_DIR}
kptgen apply pod ${KPT_BLUEPRINT_DIR} --fn-config pod-fn-config.yaml
kptgen apply webhook ${KPT_BLUEPRINT_DIR} --fn-config webhook-fn-config.yaml
kptgen apply namespace ${KPT_BLUEPRINT_DIR}
kptgen apply clusterole ${KPT_BLUEPRINT_DIR} --fn-config clusterrole-fn-config.
kptgen apply service ${KPT_BLUEPRINT_DIR} --fn-config service-fn-config.
```

TODO

```
* kptgen apply container (kube-rbac-proxy)
* kptgen apply metrics
```

## design choices

- Kind: Config is used for various resources -> webhooks, services, volumes, etc
- Namespace (root resource)
  - always created in namespace directory
  - no fnconfig
- ClusterRole kind (root resource)
  - always created in rbac directory
  - fnConfig: ClusterRole
- Pod kind (root resource)
  - created in directory based on name of the fnconfig
  - only deployment/stateful sets right now
  - permission requests: 
    - controller keyword
      - used to create clusterrole, crds are augmented here
    - other keywords are used to create roles
  - renders also service acccounts, service (optional) with deployment/statefulset
- Config kind (child resource)
  - relate always to a deployemt/statefulset (if the selector fails it fails)
  - gets crds
    - used for webhook
  - renders service
  - renders webhook -> aligned with kubebuilder
    - certificate use a specific directory with volume
  - certificate
    - follows the kubebuilder -> directory with volume
  - volume
    - if tied to certificate see above
    - if not tied to a certificate a directory is mapped
- Transaction based approach for updates
- What to do when manual changes are done
  - human change:
    - change to a rendered resource will be undone when we render again
    - if the naming changes the resources could stall resources
    - new yaml files can be added w/o issues but the human is in charge
    - Discussion: we could start from scratch every time we render, but user changes will not be possible in the blueprint
 