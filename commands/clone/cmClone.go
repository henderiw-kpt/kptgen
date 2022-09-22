package clone

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/clonedocs"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/spf13/cobra"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{
		Ctx: ctx,
	}
	c := &cobra.Command{
		Use:     "clone GIT_REPO_URL TARGET_DIR",
		Args:    cobra.MaximumNArgs(2),
		Short:   docs.CloneShort,
		Long:    docs.CloneShort + "\n" + docs.CloneLong,
		Example: docs.CloneExamples,
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

	repoURL := args[0]
	targetDir := args[1]

	// check if target directory exists, if not create it
	if err := fileutil.EnsureDir("TARGET_DIR", targetDir, true); err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "tmp")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir) // clean up

	// Clones the repository into the given dir, just as a normal git clone does
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		panic(err)
	}

	files, err := fileutil.ReadFiles(dir, true, []string{"*.yaml", "*.yml"})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		fmt.Println(file)
		rp := strings.ReplaceAll(file, dir+"/", "")
		fp := filepath.Join(targetDir, rp)

		fmt.Println(fp)
		//fmt.Println(filepath.Base(newFp))

		fileutil.EnsureDir("dummy", filepath.Dir(fp), true)
		//fmt.Println(string(yamlFile))

		//filepath.HasPrefix()
		ioutil.WriteFile(fp, yamlFile, 0644)
	}

	return nil
}
