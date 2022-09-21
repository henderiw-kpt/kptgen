package resource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	rbacv1 "k8s.io/api/rbac/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	ClusterRoleKind = "ClusterRole"
	suffixStatus    = "/status"
)

func (rn *Resource) RenderClusterRole(rules []rbacv1.PolicyRule, obj interface{}) error {
	rn.Kind = ClusterRoleKind

	// validate if crds are supplied
	// for clusterroles with crds we need to add the
	if obj != nil {
		crds, ok := obj.([]extv1.CustomResourceDefinition)
		if !ok {
			return fmt.Errorf("wrong object in renderClusterRole: %v", reflect.TypeOf(crds))
		}
		crdRulesrules := getExtraPolicyRules(crds)
		fmt.Println()
		fmt.Println(crdRulesrules)
		fmt.Println()
		rules = append(rules, crdRulesrules...)
	}
	fmt.Println(rules)

	x := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       ClusterRoleKind,
			APIVersion: rbacv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: rn.GetRoleName(),
		},
		Rules: rules,
	}

	return rn.ApplyClusterRole(x)
}

func (rn *Resource) ApplyClusterRole(x *rbacv1.ClusterRole) error {
	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	var fp string
	path, ok := x.Annotations["config.kubernetes.io/path"]
	fmt.Println(path)
	if ok {
		fp = path
		fmt.Println(rn.TargetDir)
		pathSplit := strings.Split(rn.TargetDir, "/")
		if len(pathSplit) > 1 {
			fmt.Println(pathSplit)
			pp := filepath.Join(pathSplit[:(len(pathSplit) - 1)]...)
			fmt.Println(pp)
			fp = filepath.Join(pp, fp)
		}
	}
	if fp == "" {
		fp = rn.GetFilePath(RoleSuffix)
	}

	fmt.Println(fp)

	if err := fileutil.EnsureDir(ClusterRoleKind, filepath.Dir(fp), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fp, []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}

func getExtraPolicyRules(crds []extv1.CustomResourceDefinition) []rbacv1.PolicyRule {
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
