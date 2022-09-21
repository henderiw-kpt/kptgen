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
			Name:      rn.GetName(),
			Namespace: rn.Namespace,
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
			Kind:     RoleKind,
		},
	}

	return rn.ApplyRoleBinding(x)
}

func (rn *Resource) ApplyRoleBinding(x *rbacv1.RoleBinding) error {
	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	var fp string
	path, ok := x.Annotations["config.kubernetes.io/path"]
	if ok {
		fp = path
		pathSplit := strings.Split(rn.TargetDir, "/")
		if len(pathSplit) > 1 {
			pp := filepath.Join(pathSplit[:(len(pathSplit) - 1)]...)
			fp = filepath.Join(pp, fp)
		}
	}
	if fp == "" {
		fp = rn.GetFilePath(RoleBindingSuffix)
	}

	if err := fileutil.EnsureDir(ClusterRoleKind, filepath.Dir(fp), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fp, []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}
