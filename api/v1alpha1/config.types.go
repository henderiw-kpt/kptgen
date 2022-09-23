package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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
	Group         string `json:"group,omitempty" yaml:"group,omitempty"`
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind          string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace     string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	ContainerName string `json:"containerName,omitempty" yaml:"containerName,omitempty"`
}
