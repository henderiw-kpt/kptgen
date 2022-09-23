package resource

import (
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RoleKind = "Role"
)

func (rn *Resource) RenderRole(rules []rbacv1.PolicyRule) error {
	rn.Kind = ClusterRoleKind

	x := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       RoleKind,
			APIVersion: rbacv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetRoleName(),
			Namespace: rn.GetNameSpace(),
		},
		Rules: rules,
	}

	return fileutil.CreateFileFromRObject(RoleKind, rn.GetFilePath(RoleSuffix), x)
}
