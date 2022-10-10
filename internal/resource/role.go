package resource

import (
	"fmt"
	"reflect"
	"strings"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	RoleKind = "Role"
)

func (rn *Resource) RenderRole(rules []rbacv1.PolicyRule, obj interface{}, roleName string) (*yaml.RNode, error) {

	if obj != nil && roleName == kptgenv1alpha1.ControllerClusterRoleName {
		crds, ok := obj.([]extv1.CustomResourceDefinition)
		if !ok {
			return nil, fmt.Errorf("wrong object in renderClusterRole: %v", reflect.TypeOf(crds))
		}
		crdRules := getCRDPolicyRules(crds)
		rules = append(rules, crdRules...)
	}

	x := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       RoleKind,
			APIVersion: rbacv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetRoleName(),
			Namespace: rn.GetNameSpace(),
			Labels:    rn.GetK8sLabels(),
			Annotations: map[string]string{
				kioutil.PathAnnotation:  rn.GetRelativeFilePath(RoleKind),
				kioutil.IndexAnnotation: "0",
			},
		},
		Rules: rules,
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(RoleKind), x)
}
