package apply

import (
	"context"

	"github.com/henderiw-nephio/kptgen/commands/apply/clusterrole"
	"github.com/henderiw-nephio/kptgen/commands/apply/namespace"
	"github.com/henderiw-nephio/kptgen/commands/apply/pod"
	"github.com/henderiw-nephio/kptgen/commands/apply/webhook"
	"github.com/henderiw-nephio/kptgen/commands/apply/service"
	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/applydocs"
	"github.com/spf13/cobra"
)

func GetCommand(ctx context.Context, name, version string) *cobra.Command {
	apply := &cobra.Command{
		Use:   "apply",
		Short: docs.ApplyShort,
		Long:  docs.ApplyLong,
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

	apply.AddCommand(
		pod.NewCommand(ctx, version),
		//serviceaccount.NewCommand(ctx, version),
		webhook.NewCommand(ctx, version),
		namespace.NewCommand(ctx, version),
		clusterrole.NewCommand(ctx, version),
		service.NewCommand(ctx, version),
		//container.GetCommand(ctx, version),
		//metrics.GetCommand(ctx, "", version),
	)

	return apply
}
