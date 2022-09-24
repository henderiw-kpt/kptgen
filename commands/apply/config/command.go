package config

import (
	"context"

	//kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"

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
		Use:     "config TARGET_DIR [flags]",
		Args:    cobra.MaximumNArgs(1),
		Short:   docs.ConfigShort,
		Long:    docs.ConfigShort + "\n" + docs.ConfigLong,
		Example: docs.ConfigExamples,
		RunE:    r.runE,
	}

	r.Command = c
	r.Command.Flags().StringVar(
		&r.FnConfigDir, "fn-config-dir", "", "dir where the function config files are located")
	return r
}

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command     *cobra.Command
	FnConfigDir string
	Ctx         context.Context
	pkgCfg      pkgconfig.PkgConfig
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	var err error
	r.pkgCfg, err = pkgconfig.New(args, r.FnConfigDir)
	if err != nil {
		return err
	}

	if err := r.pkgCfg.Deploy(); err != nil {
		return err
	}

	return nil
}
