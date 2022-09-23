package clusterrole

import (
	"context"

	docs "github.com/henderiw-kpt/kptgen/internal/docs/generated/applydocs"
	"github.com/henderiw-kpt/kptgen/internal/pkgconfig"
	"github.com/spf13/cobra"
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
	//TargetDir    string
	Ctx context.Context
	// dynamic input
	//pb      *kio.PackageBuffer
	//kptFile *yaml.RNode
	//fnConfig *yaml.RNode
	//fc kptgenv1alpha1.ClusterRole
	pkgCfg pkgconfig.PkgConfig
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	var err error
	r.pkgCfg, err = pkgconfig.New(args, r.FnConfigPath)
	if err != nil {
		return err
	}

	if err := r.pkgCfg.Deploy(); err != nil {
		return err
	}
	/*
		if err := r.validate(args, kptgenv1alpha1.FnClusterRoleKind); err != nil {
			return err
		}

		for name, rules := range r.fc.Spec.PermissionRequests {
			//fmt.Printf("permission requests: %s %#v\n", name, rules)
			rn := &resource.Resource{
				Operation:      resource.ClusterRoleKind,
				ControllerName: r.kptFile.GetName(),
				Name:           name,
				Namespace:      r.kptFile.GetNamespace(),
				TargetDir:      r.TargetDir,
				SubDir:         resource.RBACDir,
				NameKind:       resource.NameKindResource,
				PathNameKind:   resource.NameKindResource,
			}

			if err := rn.RenderClusterRole(rules, nil); err != nil {
				return err
			}
		}
	*/
	return nil
}

/*
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

	b, err := fileutil.ReadFile(r.FnConfigPath)
	if err != nil {
		return err
	}

	r.fc = kptgenv1alpha1.ClusterRole{}
	if err := sigyaml.Unmarshal(b, &r.fc); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	// read only Kptfile
	match := []string{"*.yaml", "*.yml", "Kptfile"}
	pb, err := pkgutil.GetPackage(r.TargetDir, match)
	if err != nil {
		return err
	}
	r.pb = pb

	cfg := config.New(r.pb, map[string]string{
		kptv1.KptFileKind: "",
		//kptgenv1alpha1.FnClusterRoleKind: filepath.Base(r.FnConfigPath),
	})

	//fmt.Println("relative", filepath.Base(r.FnConfigPath))

	selectedNodes := cfg.Get()
	if selectedNodes[kptv1.KptFileKind] == nil {
		return fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	r.kptFile = selectedNodes[kptv1.KptFileKind]

	return nil
}
*/
