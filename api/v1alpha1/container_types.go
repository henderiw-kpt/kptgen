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
	// PermissionRequests for RBAC rules required for this controller
	// to function. The RBAC manager is responsible for assessing the requested
	// permissions.
	// +optional
	PermissionRequests map[string][]rbacv1.PolicyRule `json:"permissionRequests,omitempty"`
	// Containers identifies the containers in the pod
	Containers []corev1.Container `json:"containers,omitempty"`
	// Services identifies the services the container exposes
	Service []corev1.Service `json:"services,omitempty"`
}
