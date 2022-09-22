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
	RBACDir                 = "rbac"
	ControllerSuffix        = "controller"
	NamespaceSuffix         = "namespace"
	BindingSuffix           = "binding"
	RoleSuffix              = "role"
	RoleBindingSuffix       = "role-binding"
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
	Operation      string
	NameKind       NameKind
	PathNameKind   NameKind
	ControllerName string // name of the controller/project name
	Name           string // name of the resource
	Namespace      string
	TargetDir      string
	SubDir         string
	Kind           string
}

type NameKind string

const (
	NameKindController         NameKind = "controller"
	NameKindResource           NameKind = "resource"
	NameKindControllerResource NameKind = "controllerResource"
	NameKindKind               NameKind = "kind"
	NameKindKindResource       NameKind = "kindResource"
)

func (rn *Resource) GetName() string {
	switch rn.NameKind {
	case NameKindController:
		return rn.GetControllerName("")
	case NameKindResource:
		return rn.GetResourceName("")
	case NameKindControllerResource:
		return rn.GetResourceControllerName("")
	case NameKindKind:
		return rn.GetKindName("")
	}
	return "unknown"
}

func (rn *Resource) GetFileName(extraSuffix string) string {
	switch rn.PathNameKind {
	case NameKindController:
		return rn.GetControllerName(extraSuffix)
	case NameKindResource:
		return rn.GetResourceName(extraSuffix)
	case NameKindControllerResource:
		return rn.GetResourceControllerName(extraSuffix)
	case NameKindKind:
		return rn.GetKindName(extraSuffix)
	case NameKindKindResource:
		return rn.GetKindResourceName(extraSuffix)
	}
	return "unknown"
}

func (rn *Resource) GetFilePath(extraSuffix string) string {
	return filepath.Join(rn.TargetDir, rn.SubDir, strcase.KebabCase(rn.GetFileName(extraSuffix))+".yaml")
}

func (rn *Resource) GetNameSpace() string {
	return rn.Namespace
}

func (rn *Resource) GetControllerName(extraSuffix string) string {
	if extraSuffix != "" {
		return strings.Join([]string{rn.ControllerName, ControllerSuffix, extraSuffix}, "-")
	}
	return strings.Join([]string{rn.ControllerName, ControllerSuffix}, "-")
}

func (rn *Resource) GetResourceName(extraSuffix string) string {
	if extraSuffix != "" {
		return strings.Join([]string{rn.Name, extraSuffix}, "-")
	}
	return rn.Name
}

func (rn *Resource) GetKindName(extraSuffix string) string {
	if extraSuffix != "" {
		return strings.Join([]string{rn.Kind, extraSuffix}, "-")
	}
	return rn.Kind
}

func (rn *Resource) GetKindResourceName(extraSuffix string) string {
	if extraSuffix != "" {
		return strings.Join([]string{rn.Kind, rn.Name, extraSuffix}, "-")
	}
	return strings.Join([]string{rn.Kind, rn.Name}, "-")
}

func (rn *Resource) GetResourceControllerName(extraSuffix string) string {
	if extraSuffix != "" {
		return strings.Join([]string{rn.ControllerName, rn.Name, ControllerSuffix, extraSuffix}, "-")
	}
	return strings.Join([]string{rn.ControllerName, rn.Name, ControllerSuffix}, "-")
}

func (rn *Resource) GetRoleName() string {
	return strings.Join([]string{rn.GetName(), RoleSuffix}, "-")
}

func (rn *Resource) GetRoleBindingName() string {
	return strings.Join([]string{rn.GetName(), RoleBindingSuffix}, "-")
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
	return strings.Join([]string{kptgenv1alpha1.FnConfigGroup, rn.Operation}, "/")
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
