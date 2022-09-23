package resource

import (
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RoleBindingKind = "RoleBinding"
)

func (rn *Resource) RenderRoleBinding() error {
	rn.Kind = ClusterRoleBindingKind

	x := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       RoleBindingKind,
			APIVersion: rbacv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetRoleBindingName(),
			Namespace: rn.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      rn.GetPackageResourceName(""),
				Namespace: rn.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Name:     rn.GetRoleName(),
			Kind:     RoleKind,
		},
	}

	return fileutil.CreateFileFromRObject(rn.GetFilePath(""), x)
}
