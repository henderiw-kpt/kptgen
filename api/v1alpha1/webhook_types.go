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
