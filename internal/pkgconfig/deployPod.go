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

	crds, err := resourceutil.GetCRDs(r.pb)
	if err != nil {
		return err
	}

	//fmt.Printf("crds: %#v\n", r.fc.Spec.PodTemplate)

	for roleName, rules := range fnCfg.Spec.PermissionRequests {
		rn := &resource.Resource{
			Operation:    resource.PodSuffix,
			PackageName:  r.kptFile.GetName(),
			Name:         roleName,
			Namespace:    r.kptFile.GetNamespace(),
			TargetDir:    r.targetDir,
			SubDir:       node.GetName(),
			NameKind:     resource.NameKindPackageResource,
			PathNameKind: resource.NameKindKindResource,
		}

		if roleName == kptgenv1alpha1.ControllerClusterRoleName {
			if err := rn.RenderClusterRole(rules, crds); err != nil {
				return err
			}
			if err := rn.RenderClusterRoleBinding(crds); err != nil {
				return err
			}
		} else {
			if err := rn.RenderRole(rules); err != nil {
				return err
			}
			if err := rn.RenderRoleBinding(); err != nil {
				return err
			}
		}
	}

	rn := &resource.Resource{
		Operation:    resource.PodSuffix,
		PackageName:  r.kptFile.GetName(),
		Name:         node.GetName(),
		Namespace:    r.kptFile.GetNamespace(),
		TargetDir:    r.targetDir,
		SubDir:       node.GetName(),
		NameKind:     resource.NameKindPackageResource,
		PathNameKind: resource.NameKindKind,
	}

	if err := rn.RenderServiceAccount(fnCfg.Spec); err != nil {
		return err
	}

	switch fnCfg.Spec.Type {
	case kptgenv1alpha1.DeploymentTypeDeployment:
		if err := rn.RenderDeployment(fnCfg.Spec); err != nil {
			return err
		}
	case kptgenv1alpha1.DeploymentTypeStatefulset:
		if err := rn.RenderProviderStatefulSet(fnCfg.Spec); err != nil {
			return err
		}
	}

	for _, service := range fnCfg.Spec.Services {
		if err := rn.RenderService(service, nil); err != nil {
			return err
		}
	}

	// transform the clusterrolebinding and rolebindings
	/*
		matchKind := []string{"ClusterRoleBinding", "RoleBinding"}
		for _, node := range pb.Nodes {
			for _, m := range matchKind {
				if m == node.GetKind() {
					switch node.GetKind() {
					case "ClusterRoleBinding":
						x := rbacv1.ClusterRoleBinding{}
						if err := sigyaml.Unmarshal([]byte(node.MustString()), &x); err != nil {
							return err
						}
						replace := false
						newSubjects := []rbacv1.Subject{}
						for _, subject := range x.Subjects {
							fmt.Println(subject.Kind, subject.Name, subject.Namespace)
							if subject.Kind == "ServiceAccount" {
								subject.Name = rn.GetName()
								subject.Namespace = rn.Namespace
								replace = true
							}
							newSubjects = append(newSubjects, subject)
						}
						if replace {
							x.Name = rn.GetRoleBindingName()
							x.Subjects = newSubjects
							x.Annotations = map[string]string{}
							b := new(strings.Builder)
							p := printers.YAMLPrinter{}
							if err := p.PrintObj(&x, b); err != nil {
								return err
							}

							if err := fileutil.UpdateFile(r.TargetDir, b.String(), node); err != nil {
								return err
							}

						}

					case "RoleBinding":
						x := rbacv1.RoleBinding{}
						if err := sigyaml.Unmarshal([]byte(node.MustString()), &x); err != nil {
							return err
						}
						replace := false
						newSubjects := []rbacv1.Subject{}
						for _, subject := range x.Subjects {
							fmt.Println(subject.Kind, subject.Name, subject.Namespace)
							if subject.Kind == "ServiceAccount" {
								subject.Name = rn.GetName()
								subject.Namespace = rn.Namespace
								replace = true
							}
							newSubjects = append(newSubjects, subject)
						}
						if replace {
							x.Name = rn.GetRoleBindingName()
							x.Subjects = newSubjects
							x.Annotations = map[string]string{}
							b := new(strings.Builder)
							p := printers.YAMLPrinter{}
							if err := p.PrintObj(&x, b); err != nil {
								return err
							}

							if err := fileutil.UpdateFile(r.TargetDir, b.String(), node); err != nil {
								return err
							}
						}

					}

				}
			}
		}
	*/
	return nil
}
