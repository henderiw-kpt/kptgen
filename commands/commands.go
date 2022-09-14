package commands

import (
	"context"

	"github.com/henderiw-nephio/kptgen/commands/generate"
	"github.com/spf13/cobra"
)

// GetKptCommands returns the set of kpt commands to be registered
func GetKptGenCommands(ctx context.Context, name, version string) []*cobra.Command {
	var c []*cobra.Command
	generateCmd := generate.NewCommand(ctx, name)

	c = append(c, generateCmd)
	return c
}
