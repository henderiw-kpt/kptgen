//go:generate $GOBIN/mdtogo site/reference/cli/apply internal/docs/generated/applydocs --license=none --recursive=true --strategy=cmdDocs
//go:generate $GOBIN/mdtogo site/reference/cli/clone internal/docs/generated/clonedocs --license=none --recursive=true --strategy=cmdDocs
//go:generate $GOBIN/mdtogo site/reference/cli/copy internal/docs/generated/copydocs --license=none --recursive=true --strategy=cmdDocs
//go:generate $GOBIN/mdtogo site/reference/cli/README.md internal/docs/generated/overview --license=none --strategy=cmdDocs
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/henderiw-nephio/kptgen/run"
	"github.com/spf13/cobra"
	"k8s.io/component-base/cli"
	"k8s.io/klog"
)

func main() {
	os.Exit(runMain())
}

// runMain does the initial setup in order to run kpt-gen. The return value from
// this function will be the exit code when kpt-gen terminates.
func runMain() int {
	var err error

	ctx := context.Background()

	// Enable commandline flags for klog.
	// logging will help in collecting debugging information from users
	klog.InitFlags(nil)

	cmd := run.GetMain(ctx)

	err = cli.RunNoErrOutput(cmd)
	if err != nil {
		return handleErr(cmd, err)
	}
	return 0
}

// handleErr takes care of printing an error message for a given error.
func handleErr(cmd *cobra.Command, err error) int {
	fmt.Fprintf(cmd.ErrOrStderr(), "%s \n", err.Error())
	return 1
}
