package resource

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
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

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := fileutil.EnsureDir(RoleKind, filepath.Dir(rn.GetFilePath(RoleSuffix)), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(rn.GetFilePath(RoleSuffix), []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}
