package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-immich/immich"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

const api_key_header_name = "x-api-key"

func Ptr[T any](v T) *T {
	return &v
}

func runCommand() *cli.Command {
	return &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "run the frame refresher which sync deletes and if necessary refreshes the daily images",
		Action:  runAction(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "immich-url",
				Usage:    "The URL of the immmich server including the scheme",
				Sources:  cli.EnvVars("IMMICH_URL"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "immich-api-key",
				Usage:    "The API key to connect to Immich",
				Sources:  cli.EnvVars("IMMICH_API_KEY"),
				Required: true,
			},
		},
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

		opts := []immich.ClientOption{
			immich.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
				req.Header.Set(api_key_header_name, cmd.String("immich-api-key"))

				return nil
			}),
		}

		u := cmd.String("immich-url") + "/api"

		c, err := immich.NewClientWithResponses(u, opts...)

		if err != nil {
			return err
		}

		srr, err := c.SearchRandomWithResponse(ctx, immich.SearchRandomJSONRequestBody{
			PersonIds: Ptr([]types.UUID{
				uuid.MustParse("f1a4d1c1-1ff1-4a48-93bd-21f45b23b135"),
				uuid.MustParse("796540fe-64f7-4360-8a06-7b4c1ed5ab46"),
			}),
			Type: Ptr(immich.IMAGE),
		})

		if err != nil {
			return err
		}

		if srr.HTTPResponse.StatusCode != http.StatusOK {
			return fmt.Errorf("request to get random assets failed with body :%s", string(srr.Body))
		}

		out, err := os.MkdirTemp("", "assets")
		if err != nil {
			return err
		}

		fmt.Printf("output directory is %s\n", out)

		var cnt uint16
		for _, asset := range *srr.JSON200 {
			ar, err := c.ViewAssetWithResponse(ctx, uuid.MustParse(asset.Id), &immich.ViewAssetParams{
				Size: Ptr(immich.AssetMediaSizeFullsize),
			})
			if err != nil {
				return err
			}

			f, err := os.Create(out + "/" + asset.OriginalFileName)
			if err != nil {
				return err
			}

			if _, err := f.Write(ar.Body); err != nil {
				return err
			}

			cnt++

			if cnt == 10 {
				break
			}
		}

		// TODO: ... implementation ...

		return nil
	}
}
