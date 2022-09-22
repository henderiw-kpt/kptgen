package pod

import (
	"context"
	"fmt"

	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/applydocs"
	"github.com/henderiw-nephio/kptgen/internal/resource"
	"github.com/henderiw-nephio/kptgen/internal/util/config"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/henderiw-nephio/kptgen/internal/util/pkgutil"
	"github.com/henderiw-nephio/kptgen/internal/util/resourceutil"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{
		Ctx: ctx,
	}
	c := &cobra.Command{
		Use:     "pod TARGET_DIR [flags]",
		Args:    cobra.MaximumNArgs(1),
		Short:   docs.PodShort,
		Long:    docs.PodShort + "\n" + docs.PodLong,
		Example: docs.PodExamples,
		RunE:    r.runE,
	}

	r.Command = c
	r.Command.Flags().StringVar(
		&r.FnConfigPath, "fn-config", "", "path to the function config file")
	return r
}

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command      *cobra.Command
	FnConfigPath string
	TargetDir    string
	Ctx          context.Context
	// dynamic input
	pb       *kio.PackageBuffer
	kptFile  *yaml.RNode
	fnConfig *yaml.RNode
	fc       kptgenv1alpha1.Pod
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	if err := r.validate(args, kptgenv1alpha1.FnPodKind); err != nil {
		return err
	}
	//fmt.Printf("permission requests: %#v\n", r.fc.Spec.PermissionRequests)
	//fmt.Printf("pod template: %#v\n", r.fc.Spec.PodTemplate)

	crds, err := resourceutil.GetCRDs(r.pb)
	if err != nil {
		return err
	}

	for roleName, rules := range r.fc.Spec.PermissionRequests {
		rn := &resource.Resource{
			Operation:      resource.ControllerSuffix,
			ControllerName: r.kptFile.GetName(),
			Name:           roleName,
			Namespace:      r.kptFile.GetNamespace(),
			TargetDir:      r.TargetDir,
			SubDir:         resource.RBACDir,
			NameKind:       resource.NameKindResource,
			PathNameKind:   resource.NameKindResource,
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
		Operation:      resource.ControllerSuffix,
		ControllerName: r.kptFile.GetName(),
		Name:           r.kptFile.GetName(),
		Namespace:      r.kptFile.GetNamespace(),
		TargetDir:      r.TargetDir,
		SubDir:         resource.ControllerDir,
		NameKind:       resource.NameKindController,
		PathNameKind:   resource.NameKindKind,
	}

	if err := rn.RenderServiceAccount(r.fc.Spec); err != nil {
		return err
	}

	switch r.fc.Spec.Type {
	case kptgenv1alpha1.DeploymentTypeDeployment:
		if err := rn.RenderDeployment(r.fc.Spec); err != nil {
			return err
		}
	case kptgenv1alpha1.DeploymentTypeStatefulset:
		if err := rn.RenderProviderStatefulSet(r.fc.Spec); err != nil {
			return err
		}
	}

	for _, service := range r.fc.Spec.Services {
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

func (r *Runner) validate(args []string, kind string) error {
	if len(args) < 1 {
		return fmt.Errorf("TARGET_DIR is required, positional arguments; %d provided", len(args))
	}

	r.TargetDir = args[0]

	if err := fileutil.EnsureDir("TARGET_DIR", r.TargetDir, true); err != nil {
		return err
	}

	if r.FnConfigPath == "" {
		return fmt.Errorf("a fn-config must be provided")
	}

	// read only yml, yaml files and Kptfile
	match := []string{"*.yaml", "*.yml", "Kptfile"}
	pb, err := pkgutil.GetPackage(r.TargetDir, match)
	if err != nil {
		return err
	}
	r.pb = pb

	cfg := config.New(r.pb, map[string]string{
		kptv1.KptFileKind:        "",
		kptgenv1alpha1.FnPodKind: fileutil.GetRelativePath(r.TargetDir, r.FnConfigPath),
	})

	selectedNodes := cfg.Get()
	if selectedNodes[0] == nil {
		return fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	r.kptFile = selectedNodes[0]
	if selectedNodes[1] == nil {
		return fmt.Errorf("fnConfig must be provided -> add fnConfig file with apiVersion: %s, kind: %s, name: %s", kptgenv1alpha1.FnConfigAPIVersion, kind, r.FnConfigPath)
	}
	r.fnConfig = selectedNodes[1]

	r.fc = kptgenv1alpha1.Pod{}
	if err := sigyaml.Unmarshal([]byte(r.fnConfig.MustString()), &r.fc); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}
	return nil
}
