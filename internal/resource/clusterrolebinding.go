package resource

import (
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ClusterRoleBindingKind = "ClusterRoleBinding"
)

func (rn *Resource) RenderClusterRoleBinding(obj interface{}) error {
	rn.Kind = ClusterRoleBindingKind

	x := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       ClusterRoleBindingKind,
			APIVersion: rbacv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: rn.GetName(),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      rn.GetControllerName(""),
				Namespace: rn.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Name:     rn.GetRoleName(),
			Kind:     ClusterRoleKind,
		},
	}

	return fileutil.CreateFileFromRObject(ClusterRoleBindingKind, rn.GetFilePath(RoleBindingSuffix), x)
}
