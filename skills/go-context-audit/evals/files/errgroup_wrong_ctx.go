package auditfixture

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func FetchAll(ctx context.Context, urls []string) error {
	g, gCtx := errgroup.WithContext(ctx)
	_ = gCtx

	for _, u := range urls {
		url := u
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return err
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("fetch %s: %w", url, err)
			}
			resp.Body.Close()
			return nil
		})
	}

	return g.Wait()
}
