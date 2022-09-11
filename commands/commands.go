package commands

import (
	"context"

	initialization "github.com/henderiw-nephio/kptgen/commands/init"
	"github.com/spf13/cobra"
)

// GetKptCommands returns the set of kpt commands to be registered
func GetKptGenCommands(ctx context.Context, name, version string) []*cobra.Command {
	var c []*cobra.Command
	initCmd := initialization.NewCommand(ctx, name)

	c = append(c, initCmd)
	return c
}
