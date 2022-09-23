package commands

import (
	"context"

	"github.com/henderiw-kpt/kptgen/commands/apply"
	"github.com/henderiw-kpt/kptgen/commands/clone"
	"github.com/henderiw-kpt/kptgen/commands/copy"
	"github.com/spf13/cobra"
)

// GetKptCommands returns the set of kpt commands to be registered
func GetKptGenCommands(ctx context.Context, name, version string) []*cobra.Command {
	var c []*cobra.Command
	copyCmd := copy.NewCommand(ctx, name, version)
	cloneCmd := clone.NewCommand(ctx, name, version)
	applyCmd := apply.GetCommand(ctx, name, version)

	c = append(c, copyCmd, cloneCmd, applyCmd)
	return c
}
