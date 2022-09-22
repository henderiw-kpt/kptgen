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
kpt pkg clone GIT_REPO_URL ${KPT_BLUEPRINT_DIR} // clone the fn-config templates
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
* kptgen apply service
```