package resource

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	ClusterRoleKind = "ClusterRole"
	suffixStatus    = "/status"
)

func (rn *Resource) RenderClusterRole(rules []rbacv1.PolicyRule, obj interface{}) (*yaml.RNode, error) {
	// validate if crds are supplied
	// for clusterroles with crds we need to add the
	if obj != nil {
		crds, ok := obj.([]extv1.CustomResourceDefinition)
		if !ok {
			return nil, fmt.Errorf("wrong object in renderClusterRole: %v", reflect.TypeOf(crds))
		}
		crdRules := getCRDPolicyRules(crds)
		rules = append(rules, crdRules...)
	}

	x := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       ClusterRoleKind,
			APIVersion: rbacv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: rn.GetRoleName(),
			Annotations: map[string]string{
				"config.kubernetes.io/index": "0",
				kioutil.PathAnnotation:       rn.GetRelativeFilePath(ClusterRoleKind),
				kioutil.IndexAnnotation:      "0",
			},
		},
		Rules: rules,
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(ClusterRoleKind), x)
}

func getCRDPolicyRules(crds []extv1.CustomResourceDefinition) []rbacv1.PolicyRule {
	// Our list of CRDs has no guaranteed order, so we sort them in order to
	// ensure we don't reorder our RBAC rules on each update.
	sort.Slice(crds, func(i, j int) bool { return crds[i].GetName() < crds[j].GetName() })

	groups := make([]string, 0)            // Allows deterministic iteration over groups.
	resources := make(map[string][]string) // Resources by group.
	for _, crd := range crds {
		if _, ok := resources[crd.Spec.Group]; !ok {
			resources[crd.Spec.Group] = make([]string, 0)
			groups = append(groups, crd.Spec.Group)
		}
		resources[crd.Spec.Group] = append(resources[crd.Spec.Group],
			crd.Spec.Names.Plural,
			crd.Spec.Names.Plural+suffixStatus,
		)
	}

	rules := []rbacv1.PolicyRule{}
	for _, g := range groups {
		rules = append(rules, rbacv1.PolicyRule{
			APIGroups: []string{g},
			Resources: resources[g],
			Verbs:     []string{"get", "list", "watch", "update", "patch", "create", "delete"},
		})
	}
	return rules
}
