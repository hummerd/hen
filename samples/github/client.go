package github

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/hummerd/hen"
	"github.com/hummerd/hen/middleware"
)

func NewClient() *Client {
	h := hen.NewHen(
		http.DefaultClient,
		"https://api.github.com",
		middleware.Retry(3,
			middleware.ConstBackOff(time.Second),
		),
		middleware.XRequestIDFromContext("requestID"),
		middleware.SLogger(slog.Default(), slog.LevelInfo, slog.LevelWarn, slog.LevelError),
	)

	return &Client{
		endPointGetBranches: h.NewEndpoint("/repos/%s/%s/branches", http.MethodGet,
			middleware.HTTPErrorResponse(http.StatusOK),
		),
	}
}

type Client struct {
	endPointGetBranches *hen.Endpoint
}

// GetBranches returns a list of branches for the given owner and repo.
func (c *Client) GetBranches(ctx context.Context, owner, repo string) (any, error) {
	var branches any

	err := c.endPointGetBranches.DoJSON(
		ctx,
		[]any{owner, repo},
		nil,
		nil,
		&branches,
	)

	return branches, err
}
