package v1alpha1

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	FnClusterRoleKind = "ClusterRole"
)

type ClusterRole struct {
	Spec *ClusterRoleSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ClusterRoleSpec struct {
	// PermissionRequests for RBAC rules
	// +optional
	PermissionRequests map[string][]rbacv1.PolicyRule `json:"permissionRequests,omitempty"`
}
