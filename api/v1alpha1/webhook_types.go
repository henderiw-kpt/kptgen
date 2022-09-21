package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	FnWebhookKind = "Webhook"
)

type Webhook struct {
	Spec *WebhookSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type WebhookSpec struct {
	// selector
	Selector Selector `json:"selector,omitempty" yaml:"selector,omitempty"`
	// service
	Services []corev1.Service `json:"services,omitempty"`
	// sertifcate
	Certificate Certificate `json:"certificate,omitempty" yaml:"certificate,omitempty"`
}

type Service struct {
	Port       int32 `json:"port,omitempty" yaml:"port,omitempty"`
	TargetPort int32 `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`
}

type Certificate struct {
	IssuerRef string `json:"issuerRef,omitempty" yaml:"issuerRef,omitempty"`
}

type Selector struct {
	Kind          string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	ContainerName string `json:"containerName,omitempty" yaml:"containerName,omitempty"`
}
