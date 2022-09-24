package resource

import (
	"fmt"
	"path/filepath"
	"strings"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/stoewer/go-strcase"
)

const (
	ControllerDir           = "controller"
	WebhookDir              = "webhook"
	NamespaceDir            = "namespace"
	RBACDir                 = "rbac"
	ControllerSuffix        = "controller"
	PodSuffix               = "pod"
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
	Kind string
	//NameKind     NameKind // used to define the name of the resource
	//PathNameKind NameKind // used to define the name of the path in the package
	PackageName string // name of the package/provider name
	PodName     string // name of the deployment/statefulset
	Name        string // name of the resource
	Namespace   string
	TargetDir   string
	SubDir      string
}

type NameKind string

const (
	NameKindFull            NameKind = "packagePodResource"
	NameKindPackageResource NameKind = "packageResource"
	NameKindPod             NameKind = "packagePod"
	NameKindKindResource    NameKind = "kindResource"
	NameKindKind            NameKind = "kind"
	NameKindResource        NameKind = "resource"
)

func (rn *Resource) GetNameSpace() string {
	return rn.Namespace
}

func (rn *Resource) GetResourceName(suffixes ...string) string {
	fmt.Println("GetResourceName:", rn.Kind, rn.PackageName, rn.PodName, rn.Name)
	var sb strings.Builder
	sb.WriteString(rn.PackageName)

	if rn.PodName != "" {
		sb.WriteString(fmt.Sprintf("-%s", rn.PodName))
	}
	if rn.Name != "" {
		sb.WriteString(fmt.Sprintf("-%s", rn.Name))
	}
	for _, suffix := range suffixes {
		sb.WriteString(fmt.Sprintf("-%s", suffix))
	}
	return sb.String()
}

func (rn *Resource) GetFileName(suffixes ...string) string {
	var sb strings.Builder
	sb.WriteString(strings.ToLower(rn.Kind))
	if rn.Name != "" {
		sb.WriteString(fmt.Sprintf("-%s", rn.Name))
	}
	for _, suffix := range suffixes {
		sb.WriteString(fmt.Sprintf("-%s", suffix))
	}
	return sb.String()
}

func (rn *Resource) GetFilePath(suffixes ...string) string {
	if rn.Kind == kptgenv1alpha1.FnConfigKind {
		return filepath.Join(rn.TargetDir, rn.PodName, rn.SubDir, strcase.KebabCase(rn.GetFileName(suffixes...))+".yaml")
	}
	return filepath.Join(rn.TargetDir, rn.SubDir, strcase.KebabCase(rn.GetFileName(suffixes...))+".yaml")
}

func (rn *Resource) GetServiceAccountName() string {
	return strings.Join([]string{rn.PackageName, rn.PodName}, "-")
}

func (rn *Resource) GetRoleName() string {
	return strings.Join([]string{rn.GetResourceName(), RoleSuffix}, "-")
}

func (rn *Resource) GetRoleBindingName() string {
	return strings.Join([]string{rn.GetResourceName(), RoleBindingSuffix}, "-")
}

func (rn *Resource) GetServiceName() string {
	return strings.Join([]string{rn.GetResourceName(), ServiceSuffix}, "-")
}

func (rn *Resource) GetCertificateName() string {
	return strings.Join([]string{rn.GetResourceName(), CertSuffix}, "-")
}

func (rn *Resource) GetMutatingWebhookName() string {
	return strings.Join([]string{rn.GetResourceName(), WebhookMutatingSuffix}, "-")
}

func (rn *Resource) GetValidatingWebhookName() string {
	return strings.Join([]string{rn.GetResourceName(), WebhookValidatingSuffix}, "-")
}

func (rn *Resource) GetLabelKey() string {
	if rn.Kind == ServiceSuffix {
		return strings.Join([]string{kptgenv1alpha1.FnConfigGroup, strings.Join([]string{rn.Name, rn.Kind}, "-")}, "/")
	}
	return strings.Join([]string{kptgenv1alpha1.FnConfigGroup, rn.Kind}, "/")
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
