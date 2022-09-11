package init

import (
	"context"

	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/initdocs"
	"github.com/spf13/cobra"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{
		Ctx: ctx,
	}
	c := &cobra.Command{
		Use:     "init [DIR]",
		Args:    cobra.MaximumNArgs(1),
		Short:   docs.InitShort,
		Long:    docs.InitShort + "\n" + docs.InitLong,
		Example: docs.InitExamples,
		RunE:    r.runE,
	}

	c.Flags().StringVar(&r.Description, "description", "sample description", "short description of the package.")
	c.Flags().StringSliceVar(&r.Keywords, "keywords", []string{}, "list of keywords for the package.")
	c.Flags().StringVar(&r.Site, "site", "", "link to page with information about the package.")
	//cmdutil.FixDocs("kpt", parent, c)
	r.Command = c
	return r
}

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command     *cobra.Command
	Keywords    []string
	Name        string
	Description string
	Site        string
	Ctx         context.Context
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	/*
		if len(args) == 0 {
			args = append(args, pkg.CurDir)
		}

		absPath, _, err := pathutil.ResolveAbsAndRelPaths(args[0])
		if err != nil {
			return err
		}

		pkgIniter := kptpkg.DefaultInitializer{}
		initOps := kptpkg.InitOptions{
			PkgPath:  absPath,
			RelPath:  args[0],
			Desc:     r.Description,
			Keywords: r.Keywords,
			Site:     r.Site,
		}

		return pkgIniter.Initialize(r.Ctx, filesys.FileSystemOrOnDisk{}, initOps)
	*/
	return nil
}
