package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	FnContainerKind = "Container"
)

type Container struct {
	Spec *ContainerSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ContainerSpec struct {
	// selector
	Selector Selector `json:"selector,omitempty" yaml:"selector,omitempty"`

	// ClusterRoles requested bindings
	ClusterRoles []string `json:"clusterRoles,omitempty" yaml:"clusterRoles,omitempty"`

	// PermissionRequests for RBAC rules
	// +optional
	PermissionRequests map[string][]rbacv1.PolicyRule `json:"permissionRequests,omitempty"`

	// Containers identifies the containers in the pod
	Containers []corev1.Container `json:"containers,omitempty"`

	// Services identifies the services the container exposes
	Services []corev1.Service `json:"services,omitempty"`
}
