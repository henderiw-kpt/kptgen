package resource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/yndd/ndd-runtime/pkg/utils"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	MutatingWebhookConfigurationKind   = "MutatingWebhookConfiguration"
	ValidatingWebhookConfigurationKind = "ValidatingWebhookConfiguration"
)

//crdObjects fn.KubeObjects
func (rn *Resource) RenderMutatingWebhook(cfg, obj interface{}) error {
	rn.Kind = MutatingWebhookConfigurationKind

	failurePolicy := admissionv1.Fail
	sideEffect := admissionv1.SideEffectClassNone

	webhooks := []admissionv1.MutatingWebhook{}
	objs, ok := obj.([]extv1.CustomResourceDefinition)
	if !ok {
		return fmt.Errorf("wrong object in buildValidatingWebhook: %v", reflect.TypeOf(objs))
	}
	for _, crd := range objs {
		for _, crdVersion := range crd.Spec.Versions {
			webhook := admissionv1.MutatingWebhook{
				Name:                    GetMutatingWebhookName(crd.Spec.Names.Singular, crd.Spec.Group),
				AdmissionReviewVersions: []string{"v1"},
				ClientConfig: admissionv1.WebhookClientConfig{
					Service: &admissionv1.ServiceReference{
						Name:      rn.GetServiceName(),
						Namespace: rn.Namespace,
						Path:      utils.StringPtr(strings.Join([]string{"/mutate", strings.ReplaceAll(crd.Spec.Group, ".", "-"), crdVersion.Name, crd.Spec.Names.Singular}, "-")),
					},
				},
				Rules: []admissionv1.RuleWithOperations{
					{
						Rule: admissionv1.Rule{
							APIGroups:   []string{crd.Spec.Group},
							APIVersions: []string{crdVersion.Name},
							Resources:   []string{crd.Spec.Names.Plural},
						},
						Operations: []admissionv1.OperationType{
							admissionv1.Create,
							admissionv1.Update,
						},
					},
				},
				FailurePolicy: &failurePolicy,
				SideEffects:   &sideEffect,
			}
			webhooks = append(webhooks, webhook)
		}
	}

	x := &admissionv1.MutatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       MutatingWebhookConfigurationKind,
			APIVersion: admissionv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        rn.GetMutatingWebhookName(),
			Namespace:   rn.Namespace,
			Annotations: rn.GetCertificateAnnotation(),
		},
		Webhooks: webhooks,
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := fileutil.EnsureDir(DeploymentKind, filepath.Dir(rn.GetFilePath("")), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(rn.GetFilePath(""), []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}

func (rn *Resource) RenderValidatingWebhook(cfg, obj interface{}) error {
	rn.Kind = ValidatingWebhookConfigurationKind

	failurePolicy := admissionv1.Fail
	sideEffect := admissionv1.SideEffectClassNone

	webhooks := []admissionv1.ValidatingWebhook{}
	objs, ok := obj.([]extv1.CustomResourceDefinition)
	if !ok {
		return fmt.Errorf("wrong object in buildValidatingWebhook: %v", reflect.TypeOf(objs))
	}
	for _, crd := range objs {

		//crdSingular := getCRDSingular(o)
		//crdPlural := getCRDPlural(o)
		//crdGroup := getCRDGroup(o)

		for _, crdVersion := range crd.Spec.Versions {
			webhook := admissionv1.ValidatingWebhook{
				Name:                    GetValidatingWebhookName(crd.Spec.Names.Singular, crd.Spec.Group),
				AdmissionReviewVersions: []string{"v1"},
				ClientConfig: admissionv1.WebhookClientConfig{
					Service: &admissionv1.ServiceReference{
						Name:      rn.GetServiceName(),
						Namespace: rn.Namespace,
						Path:      utils.StringPtr(strings.Join([]string{"/validate", strings.ReplaceAll(crd.Spec.Group, ".", "-"), crdVersion.Name, crd.Spec.Names.Singular}, "-")),
					},
				},
				Rules: []admissionv1.RuleWithOperations{
					{
						Rule: admissionv1.Rule{
							APIGroups:   []string{crd.Spec.Group},
							APIVersions: []string{crdVersion.Name},
							Resources:   []string{crd.Spec.Names.Plural},
						},
						Operations: []admissionv1.OperationType{
							admissionv1.Create,
							admissionv1.Update,
						},
					},
				},
				FailurePolicy: &failurePolicy,
				SideEffects:   &sideEffect,
			}
			webhooks = append(webhooks, webhook)
		}
	}

	x := &admissionv1.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       ValidatingWebhookConfigurationKind,
			APIVersion: admissionv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        rn.GetValidatingWebhookName(),
			Namespace:   rn.Namespace,
			Annotations: rn.GetCertificateAnnotation(),
		},
		Webhooks: webhooks,
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := fileutil.EnsureDir(DeploymentKind, filepath.Dir(rn.GetFilePath("")), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(rn.GetFilePath(""), []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}
