package commands

import (
	"context"

	"github.com/henderiw-nephio/kptgen/commands/add"
	"github.com/henderiw-nephio/kptgen/commands/copy"
	"github.com/spf13/cobra"
)

// GetKptCommands returns the set of kpt commands to be registered
func GetKptGenCommands(ctx context.Context, name, version string) []*cobra.Command {
	var c []*cobra.Command
	copyCmd := copy.NewCommand(ctx, name, version)
	addCmd := add.GetCommand(ctx, name, version)

	c = append(c, copyCmd, addCmd)
	return c
}
