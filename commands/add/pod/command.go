package pod

import (
	"context"
	"fmt"
	"strings"

	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/adddocs"
	"github.com/henderiw-nephio/kptgen/internal/resource"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/henderiw-nephio/kptgen/internal/util/pkgutil"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
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
	Ctx          context.Context
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("TARGET_DIR is required, positional arguments; %d provided", len(args))
	}

	targetDir := args[0]

	if err := fileutil.EnsureDir("TARGET_DIR", targetDir, true); err != nil {
		return err
	}

	if r.FnConfigPath == "" {
		return fmt.Errorf("a fn-config must be provided")
	}

	// read only yml, yaml files and Kptfile
	match := []string{"*.yaml", "*.yml", "Kptfile"}
	pb, err := pkgutil.GetPackage(targetDir, match)
	if err != nil {
		return err
	}

	kptFile, fnConfig, err := r.getConfig(pb)
	if err != nil {
		return err
	}

	fc := kptgenv1alpha1.Pod{}
	if err := sigyaml.Unmarshal([]byte(fnConfig.MustString()), &fc); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	fmt.Printf("permission requests: %#v\n", fc.Spec.PermissionRequests)

	rn := &resource.Resource{
		Suffix:    resource.ControllerSuffix,
		Name:      kptFile.GetName(),
		Namespace: kptFile.GetNamespace(),
		TargetDir: targetDir,
		SubDir:    resource.ControllerDir,
	}

	if err := resource.RenderServiceAccount(rn, fc.Spec); err != nil {
		return err
	}

	switch fc.Spec.Type {
	case kptgenv1alpha1.DeploymentTypeDeployment:
		if err := resource.RenderDeployment(rn, fc.Spec); err != nil {
			return err
		}
	case kptgenv1alpha1.DeploymentTypeStatefulset:
		if err := resource.RenderProviderStatefulSet(rn, fc.Spec); err != nil {
			return err
		}
	}

	// transform the clusterrolebinding and rolebindings
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

						if err := fileutil.UpdateFile(targetDir, b.String(), node); err != nil {
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

						if err := fileutil.UpdateFile(targetDir, b.String(), node); err != nil {
							return err
						}
					}

				}

			}
		}
	}

	return nil
}

func (r *Runner) getConfig(pb *kio.PackageBuffer) (*yaml.RNode, *yaml.RNode, error) {
	var kptFile *yaml.RNode
	var fnConfig *yaml.RNode
	for _, node := range pb.Nodes {
		if v, ok := node.GetAnnotations()[filters.LocalConfigAnnotation]; ok && v == "true" {
			if node.GetApiVersion() == kptv1.KptFileAPIVersion && node.GetKind() == kptv1.KptFileKind {
				kptFile = node
			}
			fmt.Println(node.GetName(), node.GetApiVersion(), node.GetKind())
			if node.GetApiVersion() == kptgenv1alpha1.FnConfigAPIVersion &&
				node.GetKind() == kptgenv1alpha1.FnPodKind &&
				node.GetName() == r.FnConfigPath {
				fnConfig = node
			}
		}
	}
	if kptFile == nil {
		return nil, nil, fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	if fnConfig == nil {
		return nil, nil, fmt.Errorf("fnConfig must be provided -> add fnConfig file with apiVersion: %s, kind: %s, name: %s", kptgenv1alpha1.FnConfigAPIVersion, kptgenv1alpha1.FnPodKind, r.FnConfigPath)
	}
	return kptFile, fnConfig, nil
}
