package namespace

import (
	"context"
	"fmt"

	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	docs "github.com/henderiw-kpt/kptgen/internal/docs/generated/applydocs"
	"github.com/henderiw-kpt/kptgen/internal/resource"
	"github.com/henderiw-kpt/kptgen/internal/util/config"
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	"github.com/henderiw-kpt/kptgen/internal/util/pkgutil"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
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
	Command   *cobra.Command
	TargetDir string
	Ctx       context.Context
	// dynamic input
	pb      *kio.PackageBuffer
	kptFile *yaml.RNode
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
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
