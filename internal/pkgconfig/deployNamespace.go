package pkgconfig

import (
	"github.com/henderiw-kpt/kptgen/internal/resource"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (r *pkgConfig) deployNamespace(node *yaml.RNode) error {
	rn := &resource.Resource{
		Operation:    resource.NamespaceSuffix,
		PackageName:  r.kptFile.GetName(),
		Name:         r.kptFile.GetName(),
		Namespace:    r.kptFile.GetNamespace(),
		TargetDir:    r.targetDir,
		SubDir:       resource.NamespaceDir,
		NameKind:     resource.NameKindResource,
		PathNameKind: resource.NameKindKind,
	}

	if err := resource.RenderNamespace(rn); err != nil {
		return err
	}
	return nil
}
