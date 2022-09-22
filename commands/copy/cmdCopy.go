package copy

import (
	"context"
	"fmt"

	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/copydocs"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/henderiw-nephio/kptgen/internal/util/pkgutil"
	"github.com/spf13/cobra"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{
		Ctx: ctx,
	}
	c := &cobra.Command{
		Use:     "copy SOURCE_DIR TARGET_DIR",
		Args:    cobra.MaximumNArgs(2),
		Short:   docs.CopyShort,
		Long:    docs.CopyShort + "\n" + docs.CopyLong,
		Example: docs.CopyExamples,
		RunE:    r.runE,
	}

	r.Command = c
	return r
}

func NewCommand(ctx context.Context, parent, version string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command *cobra.Command
	Ctx     context.Context
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("SOURCE_DIR and TARGET_DIR are required positional arguments; %d provided", len(args))
	}

	sourceDir := args[0]
	targetDir := args[1]

	// check if src directory exists but dont create it
	if err := fileutil.EnsureDir("SOURCE_DIR", sourceDir, false); err != nil {
		return err
	}

	// check if target directory exists, if not create it
	if err := fileutil.EnsureDir("TARGET_DIR", targetDir, true); err != nil {
		return err
	}

	// read only yml, yaml files
	m := []string{"*.yaml", "*.yml"}
	pb, err := pkgutil.GetPackage(sourceDir, m)
	if err != nil {
		return err
	}

	// the files with the following kinds and names will be copied
	matchKinds := map[string][]string{
		//"ClusterRole":              {"manager-role"},
		//"ClusterRoleBinding":       {"manager-rolebinding"},
		//"Role":                     {"leader-election-role"},
		//"RoleBinding":              {"leader-election-rolebinding"},
		"CustomResourceDefinition": {""},
	}

	for _, node := range pb.Nodes {
		if names, ok := matchKinds[node.GetKind()]; ok {
			for _, name := range names {
				if name == "" || node.GetName() == name {
					fileutil.CreateFileFromRNode(targetDir, node)
				}
			}
		}
	}
	return nil
}
