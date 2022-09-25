package pkgconfig

import (
	"github.com/henderiw-kpt/kptgen/internal/resource"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (r *pkgConfig) deployNamespace(n *yaml.RNode) error {
	// render the naemspace based on the fnConfig input
	rn := &resource.Resource{
		Kind:             resource.NamespaceSuffix,
		PackageName:      r.kptFile.GetName(),
		Namespace:        r.kptFile.GetNamespace(),
		TargetDir:        r.targetDir,
		SubDir:           resource.NamespaceDir,
		ResourceNameKind: resource.NameKindKind,
	}

	node, err := rn.RenderNamespace()
	if err != nil {
		return err
	}
	r.pkgResources.Add(node)
	return nil
}
