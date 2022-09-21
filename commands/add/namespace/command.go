package namespace

import (
	"context"
	"fmt"

	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/adddocs"
	"github.com/henderiw-nephio/kptgen/internal/resource"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/henderiw-nephio/kptgen/internal/util/pkgutil"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
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
	Command *cobra.Command
	Ctx     context.Context
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("TARGET_DIR is required, positional arguments; %d provided", len(args))
	}

	targetDir := args[0]

	if err := fileutil.EnsureDir("TARGET_DIR", targetDir, true); err != nil {
		return err
	}

	// read only yml, yaml files and Kptfile
	match := []string{"Kptfile"}
	pb, err := pkgutil.GetPackage(targetDir, match)
	if err != nil {
		return err
	}

	kptFile, err := r.getConfig(pb)
	if err != nil {
		return err
	}

	rn := &resource.Resource{
		Suffix:    resource.NamespaceSuffix,
		Name:      kptFile.GetName(),
		Namespace: kptFile.GetNamespace(),
		TargetDir: targetDir,
		SubDir:    resource.NamespaceDir,
	}

	if err := resource.RenderNamespace(rn); err != nil {
		return err
	}

	return nil
}

func (r *Runner) getConfig(pb *kio.PackageBuffer) (*yaml.RNode, error) {
	var kptFile *yaml.RNode
	for _, node := range pb.Nodes {
		if v, ok := node.GetAnnotations()[filters.LocalConfigAnnotation]; ok && v == "true" {
			if node.GetApiVersion() == kptv1.KptFileAPIVersion && node.GetKind() == kptv1.KptFileKind {
				kptFile = node
			}
		}
	}
	if kptFile == nil {
		return nil, fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	return kptFile, nil
}
