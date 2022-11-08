# kptgen

kptgen is a tool to help build a kpt package from a kubebuilder project. kptgen uses function config files to determine which application manifests need to be rendered in the package. function configs are modelled as KRM resources (see examples in the fn-config section)

## usage

install kptgen

```
go install -v github.com/henderiw-kpt/kptgen@$v0.0.9
```

Create the directories where the package resides. Right now kptgen seperates the CRD and the app specific manifests in different sub-packages. On top kptgen uses a function config files files to determine which application manifests need to be rendered in the package. We recommend to put these function config files in a seperate directory from the real application package.

```
PROJECT ?= <your project name - typically the controller name>
KPT_BLUEPRINT_CFG_DIR = blueprint/fn-config
KPT_BLUEPRINT_PKG_DIR = blueprint/${PROJECT}

mkdir -p ${KPT_BLUEPRINT_CFG_DIR}
mkdir -p ${KPT_BLUEPRINT_PKG_DIR}/crd/bases
mkdir -p ${KPT_BLUEPRINT_PKG_DIR}/app
```

For namespaces we use the namespace in the Kptfile

```
kpt pkg init ${KPT_BLUEPRINT_PKG_DIR} --description "${PROJECT} controller"
```

add the namespace in the generated kptfile

```
apiVersion: kpt.dev/v1
kind: Kptfile
metadata:
  name: admin
  namespace: ndd-system
  annotations:
    config.kubernetes.io/local-config: "true"
info:
  description: admin controller
```

Add the function config files in the fn-config directory we created in step 1. Normally the pod fn-config is sufficient.
See example in the fn-config where it is also explained the changes one needs to do to the example


One this is prepared we can create the CRD sub-package and render through controlelr-gen the CRDs from kubebuilder in the package.

```
kpt pkg init ${KPT_BLUEPRINT_PKG_DIR}/crd --description "${PROJECT} crd"
controller-gencrd paths="./..." output:crd:artifacts:config=${KPT_BLUEPRINT_PKG_DIR}/crd/bases
```

After we can render the aplication package. The order matters since the CRDs are used for the permissions in the controller role/cluster-role

```
kpt pkg init ${KPT_BLUEPRINT_PKG_DIR}/app --description "${PROJECT} app"
kptgen apply config ${KPT_BLUEPRINT_PKG_DIR} --fn-config-dir ${KPT_BLUEPRINT_CFG_DIR}
```

to deal with some package-context overlap in sub-packages we remove the package-context.yam files

```
rm ${KPT_BLUEPRINT_PKG_DIR}/package-context.yaml
rm ${KPT_BLUEPRINT_PKG_DIR}/crd/package-context.yaml
rm ${KPT_BLUEPRINT_PKG_DIR}/app/package-context.yaml
```

One can automate this through a makefile, which is shown below

## example makefile extensions

Makefile extension to leverage kptgen

```
VERSION ?= latest
REGISTRY ?= yndd
PROJECT ?= admin

KPT_BLUEPRINT_CFG_DIR ?= blueprint/fn-config
KPT_BLUEPRINT_PKG_DIR ?= blueprint/${PROJECT}

.PHONY: generate
generate: controller-gen kpt kptgen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	mkdir -p ${KPT_BLUEPRINT_CFG_DIR}
	mkdir -p ${KPT_BLUEPRINT_PKG_DIR}/crd/bases
	mkdir -p ${KPT_BLUEPRINT_PKG_DIR}/app
	$(CONTROLLER_GEN) crd paths="./..." output:crd:artifacts:config=${KPT_BLUEPRINT_PKG_DIR}/crd/bases
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	kpt pkg init ${KPT_BLUEPRINT_PKG_DIR} --description "${PROJECT} controller"
	kpt pkg init ${KPT_BLUEPRINT_PKG_DIR}/crd --description "${PROJECT} crd"
	kpt pkg init ${KPT_BLUEPRINT_PKG_DIR}/app --description "${PROJECT} app"
	kptgen apply config ${KPT_BLUEPRINT_PKG_DIR} --fn-config-dir ${KPT_BLUEPRINT_CFG_DIR}
	rm ${KPT_BLUEPRINT_PKG_DIR}/package-context.yaml
	rm ${KPT_BLUEPRINT_PKG_DIR}/crd/package-context.yaml
	rm ${KPT_BLUEPRINT_PKG_DIR}/app/package-context.yaml
```

to use the right kptgen executables we can install them locally in the project directory 

```
KPT ?= $(LOCALBIN)/kpt
KPTGEN ?= $(LOCALBIN)/kptgen

KPT_VERSION ?= main
KPTGEN_VERSION ?= v0.0.9

.PHONY: kpt
kpt: $(KPT) ## Download kpt locally if necessary.
$(KPT): $(LOCALBIN)
	test -s $(LOCALBIN)/kpt || GOBIN=$(LOCALBIN) go install -v github.com/GoogleContainerTools/kpt@$(KPT_VERSION)

.PHONY: kptgen
kptgen: $(KPTGEN) ## Download kptgen locally if necessary.
$(KPTGEN): $(LOCALBIN)
	test -s $(LOCALBIN)/kptgen || GOBIN=$(LOCALBIN) go install -v github.com/henderiw-kpt/kptgen@$(KPTGEN_VERSION)
```

## function-config

There are x function-config options. We can render them individually

```
kptgen apply pod -fn-config <pod-fn-cfg-file>
```

or kptgen can render all configuration from a directory that contains all the fn-configs for the packages

```
kptgen apply config --fn-config-dir <fn-config-directory>
```

### POD function-config

The main function config is Pod/fn.kptgen.dev/v1alpha1, which is used to render the manifest for a k8s controller. It geenrates the following manifests based on the fn-config file:
- deployment or statefulset
- clusterrole/clusterrole bindings from permission setting in the fn-config
- role/role bindings based on permission setting in the fn-config
- service-account

Typically one configures the function config from the example:
- type: deployment | statefulset
- permissions:
  - crds are automatically taken care of, but the user can specify to use them as scope: role or cluster. This is done by setting the scope to either cluster or role within the permission request.
  - Any additional permissions beyond the CRDs need to be explicitaly mentioned in the fn-config as per example below.
- image need to be set to the one you use

```
apiVersion: fn.kptgen.dev/v1alpha1
kind: Pod
metadata:
  name: controller
  annotations:
    config.kubernetes.io/local-config: "true"
  namespace: ndd-system
spec:
  type: deployment
  replicas: 1
  permissionRequests:
    controller:
      scope: cluster
      permissions:
      - apiGroups: ["*"]
        resources: [events]
        verbs: [get, list, watch, update, patch, create, delete]
    porch:
      scope: cluster
      permissions:
      - apiGroups: [porch.kpt.dev]
        resources: [packagerevisionresources, packagerevisions]
        verbs: [get, list, watch, update, patch, create, delete]
      - apiGroups: [porch.kpt.dev]
        resources: [packagerevisionresources/status, packagerevisions/status, packagerevisions/approval]
        verbs: [get, update, patch]
    leader-election:
      permissions:
      - apiGroups: [""]
        resources: [configmaps]
        verbs: [get, list, watch, update, patch, create, delete]
      - apiGroups: [coordination.k8s.io]
        resources: [leases]
        verbs: [get, list, watch, update, patch, create, delete]
      - apiGroups: [""]
        resources: [events]
        verbs: [create, patch]
  template:
    spec:
      containers:
      - name: kube-rbac-proxy
        image: gcr.io/kubebuilder/kube-rbac-proxy:v0.8.0
        args:
        - --secure-listen-address=0.0.0.0:8443
        - --upstream=http://127.0.0.1:8080/
        - --logtostderr=true
        - --v=10
        ports:
        - containerPort: 8443
          name: https
      - name: controller
        image: yndd/admin-controller:latest
        command:
        - /manager
        args:
        - --health-probe-bind-address=:8081
        - --metrics-bind-address=127.0.0.1:8080
        - --leader-elect
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
        # TODO(user): Configure the resources accordingly based on the project requirements.
        # More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 10m
            memory: 64Mi
        env:
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
```

### Service function-config

Additional services can be modelled through a specific service config. Example of a metric service that is attached to the deployment.

The selector refers to the main deployment such that we augment the deployment manifest with the proper information labels. The services section contains the specific config for the service that is consumed.

```
apiVersion: fn.kptgen.dev/v1alpha1
kind: Config
metadata:
  name: metrics
  annotations:
    config.kubernetes.io/local-config: "true"
  namespace: ndd-system
spec:
  selector:
    kind: Deployment
    name: controller
    containerName: controller
  services:
  - spec:
      ports:
      - name: metrics
        port: 8443
        targetPort: 443
        protocol: TCP
```

### webhook function-config

Allows to render the additional manifests needed to add a webhook to the deployment package

```
apiVersion: fn.kptgen.dev/v1alpha1
kind: Webhook
metadata:
  name: webhook
  annotations:
    config.kubernetes.io/local-config: "true"
  namespace: ndd-system
spec:
  selector:
    kind: Deployment
    name: controller
    containerName: controller
  services:
  - spec:
      ports:
      - name: webhook
        port: 9443
        targetPort: 9443
        protocol: TCP
  certificate: self-signed
```

### namespace is optional and can be modelled seperately

since the namespace can be optional it is not added by default and can be explicitly set using 


## example

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
      - used to augment crds to
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
 