package repo

import (
  "context"
  "time"
)

const defaultTimeout = 5 * time.Second

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
  if ctx == nil {
    ctx = context.Background()
  }
  return context.WithTimeout(ctx, defaultTimeout)
}
