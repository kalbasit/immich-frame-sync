package cmd

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func runCommand() *cli.Command {
	return &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "run the frame refresher which sync deletes and if necessary refreshes the daily images",
		Action:  runAction(),
		Flags:   []cli.Flag{},
	}
}

func runAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		logger := zerolog.Ctx(ctx).With().Str("cmd", "run").Logger()

		ctx = logger.WithContext(ctx)

		ctx, cancel := context.WithCancel(ctx)

		g, ctx := errgroup.WithContext(ctx)
		defer func() {
			if err := g.Wait(); err != nil {
				logger.Error().Err(err).Msg("error returned from g.Wait()")
			}
		}()

		// NOTE: Reminder that defer statements run last to first so the first
		// thing that happens here is the context is canceled which triggers the
		// errgroup 'g' to start exiting.
		defer cancel()

		g.Go(func() error {
			return autoMaxProcs(ctx, 30*time.Second, logger)
		})

		// TODO: ... implementation ...

		return nil
	}
}
