package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

const (
	FnConfigKind = "Config"
)

type Config struct {
	Spec *ConfigSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ConfigSpec struct {
	// selector
	Selector Selector `json:"selector,omitempty" yaml:"selector,omitempty"`

	// webhook
	Webhook bool `json:"webhook,omitempty" yaml:"webhook,omitempty"`

	// volume
	Volume bool `json:"volume,omitempty" yaml:"volume,omitempty"`

	// service
	Services []corev1.Service `json:"services,omitempty" yaml:"services,omitempty"`

	// sertifcate
	Certificate Certificate `json:"certificate,omitempty" yaml:"certificate,omitempty"`

	// ClusterRoles requested bindings
	ClusterRoles []string `json:"clusterRoles,omitempty" yaml:"clusterRoles,omitempty"`

	// PermissionRequests for RBAC rules
	// +optional
	PermissionRequests map[string][]rbacv1.PolicyRule `json:"permissionRequests,omitempty"`

	// Containers identifies the containers in the pod
	Containers []corev1.Container `json:"containers,omitempty"`
}

type Certificate struct {
	IssuerRef string `json:"issuerRef,omitempty" yaml:"issuerRef,omitempty"`
}

type Selector struct {
	// ResId refers to a GVKN/Ns of a resource.
	resid.ResId `json:",inline,omitempty" yaml:",inline,omitempty"`

	ContainerName string `json:"containerName,omitempty" yaml:"containerName,omitempty"`
}
