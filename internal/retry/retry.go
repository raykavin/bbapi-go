package retry

import (
	"context"
	"time"

	gkretry "github.com/raykavin/gokit/retry"
)

// Do wraps the currently available retry dependency behind an internal package
// so the transport can evolve without leaking third-party retry details.
func Do(
	ctx context.Context,
	maxAttempts int,
	waitMin time.Duration,
	waitMax time.Duration,
	shouldRetry func(attempt int, err error) bool,
	fn func() error,
) error {
	return gkretry.Do(
		ctx,
		maxAttempts,
		waitMin,
		waitMax,
		shouldRetry,
		fn,
	)
}
