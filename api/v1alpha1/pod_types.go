package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	ControllerClusterRoleName = "controller"
)

const (
	FnPodKind = "Pod"
)

type DeploymentType string

const (
	DeploymentTypeStatefulset DeploymentType = "statefulset"
	DeploymentTypeDeployment  DeploymentType = "deployment"
)

type PolicyScope string

const (
	PolicyScopeCluster   PolicyScope = "cluster"
	PolicyScopeNamespace PolicyScope = "namespace"
)

type Pod struct {
	Spec *PodSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type PodSpec struct {
	// Type is the type of the deployment
	// +kubebuilder:validation:Enum=`statefulset`;`deployment`
	// +kubebuilder:default=deployment
	Type DeploymentType `json:"type,omitempty"`

	// Replicas defines the amount of replicas expected
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`

	// MaxReplicas defines the max expected replications of this pod
	// +kubebuilder:default=8
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`

	// MaxJobNumber indication on how many jobs a given pods should hold
	MaxJobNumber *int32 `json:"maxJobNumber,omitempty"`

	// ClusterRoles requested bindings
	ClusterRoles []string `json:"clusterRoles,omitempty" yaml:"clusterRoles,omitempty"`

	// PermissionRequests for RBAC rules required for this controller to function.
	// +optional
	PermissionRequests map[string]*PolicyRules `json:"permissionRequests,omitempty"`

	// pods define the pod specification used by the controller for LCM/resource allocation
	PodTemplate corev1.PodTemplateSpec `json:"template,omitempty"`
	// Services identifies the services the pod exposes
	Services []corev1.Service `json:"services,omitempty"`
}

type PolicyRules struct {
	// Scope defines the scope of the policy rules
	// +kubebuilder:validation:Enum=`namespace`;`cluster`
	// +kubebuilder:default=namespace
	Scope PolicyScope `json:"scope,omitempty"`
	// rules is the set of rules the permissions requests should fullfil
	Permissions []rbacv1.PolicyRule `json:"permissions"`
}
