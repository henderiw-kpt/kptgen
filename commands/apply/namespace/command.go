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
	Ctx     context.Context
	pkgCfg  pkgconfig.PkgConfig
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

	return nil
}
