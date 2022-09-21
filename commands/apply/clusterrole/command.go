package clusterrole

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/applydocs"
	"github.com/henderiw-nephio/kptgen/internal/resource"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/henderiw-nephio/kptgen/internal/util/pkgutil"
	"github.com/spf13/cobra"
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
		Use:     "clusterrole TARGET_DIR [flags]",
		Args:    cobra.MaximumNArgs(1),
		Short:   docs.ClusterroleShort,
		Long:    docs.ClusterroleShort + "\n" + docs.ClusterroleLong,
		Example: docs.ClusterroleExamples,
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
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
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

	kptFile, fnConfig, err := r.getConfig(pb)
	if err != nil {
		return err
	}

	fc := kptgenv1alpha1.ClusterRole{}
	if err := sigyaml.Unmarshal([]byte(fnConfig.MustString()), &fc); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	for name, rules := range fc.Spec.PermissionRequests {
		//fmt.Printf("permission requests: %s %#v\n", name, rules)

		rn := &resource.Resource{
			Operation:      resource.ClusterRoleKind,
			ControllerName: kptFile.GetName(),
			Name:           name,
			Namespace:      kptFile.GetNamespace(),
			TargetDir:      r.TargetDir,
			SubDir:         resource.RBACDir,
			NameKind:       resource.NameKindResource,
			PathNameKind:   resource.NameKindResource,
		}

		fmt.Println(rn)

		if err := rn.RenderClusterRole(rules, nil); err != nil {
			return err
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
				node.GetKind() == kptgenv1alpha1.FnClusterRoleKind &&
				getResosurcePathFromConfigPath(r.TargetDir, r.FnConfigPath) == node.GetAnnotations()["internal.config.kubernetes.io/path"] {
				fnConfig = node
			}
		}
	}
	if kptFile == nil {
		return nil, nil, fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	if fnConfig == nil {
		return nil, nil, fmt.Errorf("fnConfig must be provided -> add fnConfig file with apiVersion: %s, kind: %s, name: %s", kptgenv1alpha1.FnConfigAPIVersion, kptgenv1alpha1.FnClusterRoleKind, r.FnConfigPath)
	}
	return kptFile, fnConfig, nil
}

func getResosurcePathFromConfigPath(targetDir, configPath string) string {
	split1 := strings.Split(targetDir, "/")
	split2 := strings.Split(configPath, "/")

	idx := 0
	for i := range split1 {
		if split1[i] != split2[i] {
			break
		}
		idx = i
	}
	if idx > 0 {
		configPath = filepath.Join(split2[(idx):]...)
	}
	return configPath
}
