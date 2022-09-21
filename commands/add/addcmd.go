package add

import (
	"context"

	"github.com/henderiw-nephio/kptgen/commands/add/namespace"
	"github.com/henderiw-nephio/kptgen/commands/add/pod"
	"github.com/henderiw-nephio/kptgen/commands/add/webhook"
	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/adddocs"
	"github.com/spf13/cobra"
)

func GetCommand(ctx context.Context, name, version string) *cobra.Command {
	add := &cobra.Command{
		Use:   "add",
		Short: docs.AddShort,
		Long:  docs.AddLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := cmd.Flags().GetBool("help")
			if err != nil {
				return err
			}
			if h {
				return cmd.Help()
			}
			return cmd.Usage()
		},
	}

	add.AddCommand(
		pod.NewCommand(ctx, version),
		//serviceaccount.NewCommand(ctx, version),
		webhook.NewCommand(ctx, version),
		namespace.NewCommand(ctx, version),
		//container.GetCommand(ctx, version),
		//metrics.GetCommand(ctx, "", version),
	)

	return add
}
