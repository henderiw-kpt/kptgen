package resource

import (
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	RoleBindingKind = "RoleBinding"
)

func (rn *Resource) RenderRoleBinding() (*yaml.RNode, error) {
	x := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       RoleBindingKind,
			APIVersion: rbacv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetRoleBindingName(),
			Namespace: rn.Namespace,
			Annotations: map[string]string{
				"config.kubernetes.io/index": "0",
				kioutil.PathAnnotation:       rn.GetRelativeFilePath(RoleBindingKind),
				kioutil.IndexAnnotation:      "0",
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      rn.GetServiceAccountName(),
				Namespace: rn.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Name:     rn.GetRoleName(),
			Kind:     RoleKind,
		},
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(RoleBindingKind), x)
}
