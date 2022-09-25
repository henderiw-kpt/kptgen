package pkgconfig

import (
	"fmt"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/resource"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

func (r *pkgConfig) deployClusterRole(node *yaml.RNode) error {
	// marshal the fnConfig
	fnCfg := kptgenv1alpha1.ClusterRole{}
	if err := sigyaml.Unmarshal([]byte(node.MustString()), &fnCfg); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	// render the cluster roles based on the fnConfig input
	for roleName, rules := range fnCfg.Spec.PermissionRequests {
		//fmt.Printf("permission requests: %s %#v\n", name, rules)
		rn := &resource.Resource{
			Kind:             kptgenv1alpha1.FnClusterRoleKind,
			PackageName:      r.kptFile.GetName(),
			Name:             roleName,
			Namespace:        r.kptFile.GetNamespace(),
			TargetDir:        r.targetDir,
			SubDir:           resource.RBACDir,
			ResourceNameKind: resource.NameKindKindResource,
		}

		node, err := rn.RenderClusterRole(rules, nil)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
	}
	return nil
}
