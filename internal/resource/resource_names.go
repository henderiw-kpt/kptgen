package resource

import (
	"path/filepath"
	"strings"

	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	"github.com/stoewer/go-strcase"
)

const (
	ControllerDir           = "controller"
	WebhookDir              = "webhook"
	NamespaceDir            = "namespace"
	ControllerSuffix        = "controller"
	NamespaceSuffix         = "namespace"
	RoleBindingSuffix       = "role-binding"
	RoleSuffix              = "role"
	WebhookSuffix           = "webhook"
	ServiceSuffix           = "svc"
	CertSuffix              = "serving-cert"
	CertPathSuffix          = "serving-certs"
	CertInjectionKey        = "cert-manager.io/inject-ca-from"
	WebhookMutatingSuffix   = "mutating-configuration"
	WebhookValidatingSuffix = "validating-configuration"
)

type Resource struct {
	//Prefix    string
	Suffix    string
	Name      string
	Namespace string
	TargetDir string
	SubDir    string
	Kind      string
}

func (rn *Resource) GetFilePath() string {
	return filepath.Join(rn.TargetDir, rn.SubDir, strcase.KebabCase(rn.Kind)+".yaml")
}

func (rn *Resource) GetName() string {
	return strings.Join([]string{rn.Name, rn.Suffix}, "-")
}

func (rn *Resource) GetRoleBindingName() string {
	return strings.Join([]string{rn.GetName(), RoleBindingSuffix}, "-")
}

func (rn *Resource) GetRoleName() string {
	return strings.Join([]string{rn.GetName(), RoleSuffix}, "-")
}

func (rn *Resource) GetServiceName() string {
	return strings.Join([]string{rn.GetName(), ServiceSuffix}, "-")
}

func (rn *Resource) GetCertificateName() string {
	return strings.Join([]string{rn.GetName(), CertSuffix}, "-")
}

func (rn *Resource) GetMutatingWebhookName() string {
	return strings.Join([]string{rn.GetName(), WebhookMutatingSuffix}, "-")
}

func (rn *Resource) GetValidatingWebhookName() string {
	return strings.Join([]string{rn.GetName(), WebhookValidatingSuffix}, "-")
}

func (rn *Resource) GetLabelKey() string {
	return strings.Join([]string{kptgenv1alpha1.FnConfigGroup, rn.Suffix}, "/")
}

func (rn *Resource) GetDnsName(x ...string) string {
	s := []string{rn.GetServiceName(), rn.Namespace, ServiceSuffix}
	if len(x) > 0 {
		s = append(s, x...)
	}
	return strings.Join(s, ".")
}

func (rn *Resource) GetCertificateAnnotation() map[string]string {
	return map[string]string{
		CertInjectionKey: strings.Join([]string{rn.Namespace, rn.GetCertificateName()}, "/"),
	}
}

func GetMutatingWebhookName(crdSingular, crdGroup string) string {
	return strings.Join([]string{"m" + crdSingular, crdGroup}, ".")
}

func GetValidatingWebhookName(crdSingular, crdGroup string) string {
	return strings.Join([]string{"v" + crdSingular, crdGroup}, ".")
}
