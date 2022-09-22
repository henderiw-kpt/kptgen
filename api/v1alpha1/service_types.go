package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	FnServiceKind = "Service"
)

type Service struct {
	Spec *ServiceSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ServiceSpec struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// selector
	Selector Selector `json:"selector,omitempty" yaml:"selector,omitempty"`
	// service
	Services []corev1.Service `json:"services,omitempty"`
}
