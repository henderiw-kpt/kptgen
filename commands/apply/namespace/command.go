package namespace

import (
	"context"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
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
		Use:     "namespace TARGET_DIR [flags]",
		Args:    cobra.MaximumNArgs(1),
		Short:   docs.NamespaceShort,
		Long:    docs.NamespaceShort + "\n" + docs.NamespaceLong,
		Example: docs.NamespaceExamples,
		RunE:    r.runE,
	}

	r.Command = c
	return r
}

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command *cobra.Command
	//TargetDir string
	Ctx context.Context
	// dynamic input
	//pb      *kio.PackageBuffer
	//kptFile *yaml.RNode
	pkgCfg pkgconfig.PkgConfig
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	var err error
	// namespace is used as a dummyFnCOnfig
	r.pkgCfg, err = pkgconfig.New(args, kptgenv1alpha1.DummyFnConfig)
	if err != nil {
		return err
	}

	if err := r.pkgCfg.Deploy(); err != nil {
		return err
	}

	/*
		if err := r.validate(args, "Namespace"); err != nil {
			return err
		}

		rn := &resource.Resource{
			Operation:      resource.NamespaceSuffix,
			ControllerName: r.kptFile.GetName(),
			Name:           r.kptFile.GetName(),
			Namespace:      r.kptFile.GetNamespace(),
			TargetDir:      r.TargetDir,
			SubDir:         resource.NamespaceDir,
			NameKind:       resource.NameKindResource,
			PathNameKind:   resource.NameKindKind,
		}

		if err := resource.RenderNamespace(rn); err != nil {
			return err
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

	if err := fileutil.EnsureDir(r.TargetDir, true); err != nil {
		return err
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
	})

	selectedNodes := cfg.Get()
	if selectedNodes[kptv1.KptFileKind] == nil {
		return fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	r.kptFile = selectedNodes[kptv1.KptFileKind]

	return nil
}
*/
