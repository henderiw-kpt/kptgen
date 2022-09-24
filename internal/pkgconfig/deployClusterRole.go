package pkgconfig

import (
	"fmt"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/resource"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

func (r *pkgConfig) deployClusterRole(node *yaml.RNode) error {
	fnCfg := kptgenv1alpha1.ClusterRole{}
	if err := sigyaml.Unmarshal([]byte(node.MustString()), &fnCfg); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	for roleName, rules := range fnCfg.Spec.PermissionRequests {
		//fmt.Printf("permission requests: %s %#v\n", name, rules)
		rn := &resource.Resource{
			Kind:        kptgenv1alpha1.FnClusterRoleKind,
			PackageName: r.kptFile.GetName(),
			Name:        roleName,
			Namespace:   r.kptFile.GetNamespace(),
			TargetDir:   r.targetDir,
			SubDir:      resource.RBACDir,
			//NameKind:     resource.NameKindPackageResource,
			//PathNameKind: resource.NameKindKindResource,
		}

		if err := rn.RenderClusterRole(rules, nil); err != nil {
			return err
		}
	}
	return nil
}
