package pkgconfig

import (
	"fmt"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/resource"
	"github.com/henderiw-kpt/kptgen/internal/util/resourceutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

func (r *pkgConfig) deployPod(node *yaml.RNode) error {
	fnCfg := kptgenv1alpha1.Pod{}
	if err := sigyaml.Unmarshal([]byte(node.MustString()), &fnCfg); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	//fmt.Printf("permission requests: %#v\n", r.fc.Spec.PermissionRequests)
	//fmt.Printf("pod template: %#v\n", r.fc.Spec.PodTemplate)

	crds, err := resourceutil.GetCRDs(r.pkgResources)
	if err != nil {
		return err
	}

	//fmt.Printf("crds: %#v\n", r.fc.Spec.PodTemplate)

	for roleName, rules := range fnCfg.Spec.PermissionRequests {
		rn := &resource.Resource{
			Kind:             kptgenv1alpha1.FnPodKind,
			PackageName:      r.kptFile.GetName(),
			PodName:          node.GetName(),
			Name:             roleName,
			Namespace:        r.kptFile.GetNamespace(),
			TargetDir:        r.targetDir,
			SubDir:           node.GetName(),
			ResourceNameKind: resource.NameKindFull,
			//PathNameKind: resource.NameKindKindResource,
		}

		if rules.Scope == kptgenv1alpha1.PolicyScopeCluster {
			node, err := rn.RenderClusterRole(rules.Permissions, crds, roleName)
			if err != nil {
				return err
			}
			r.pkgResources.Add(node)
			node, err = rn.RenderClusterRoleBinding()
			if err != nil {
				return err
			}
			r.pkgResources.Add(node)
		} else {
			node, err := rn.RenderRole(rules.Permissions, crds, roleName)
			if err != nil {
				return err
			}
			r.pkgResources.Add(node)
			node, err = rn.RenderRoleBinding()
			if err != nil {
				return err
			}
			r.pkgResources.Add(node)
		}

	}

	rn := &resource.Resource{
		Kind:             kptgenv1alpha1.FnPodKind,
		PackageName:      r.kptFile.GetName(),
		PodName:          node.GetName(),
		Name:             node.GetName(),
		Namespace:        r.kptFile.GetNamespace(),
		TargetDir:        r.targetDir,
		SubDir:           node.GetName(),
		ResourceNameKind: resource.NameKindPackagePod,
		//PathNameKind: resource.NameKindKindResource,
	}

	node, err = rn.RenderServiceAccount(fnCfg.Spec)
	if err != nil {
		return err
	}
	r.pkgResources.Add(node)

	switch fnCfg.Spec.Type {
	case kptgenv1alpha1.DeploymentTypeDeployment:
		node, err := rn.RenderDeployment(fnCfg.Spec)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
	case kptgenv1alpha1.DeploymentTypeStatefulset:
		node, err := rn.RenderProviderStatefulSet(fnCfg.Spec)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
	}

	for _, service := range fnCfg.Spec.Services {
		node, err := rn.RenderService(service, nil)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
	}

	return nil
}
