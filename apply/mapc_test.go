package apply_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestMapC(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.MapC(2, func(_ context.Context, s string) (string, error) {
			return s + s, nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"testtest", "foofoo", "bazzbazz", "traintrain"})
	})

	t.Run("runs concurrently", func(t *testing.T) {
		t.Parallel()

		start := time.Now()

		delay := 200 * time.Millisecond

		out := testhelpers.TestStage(t, apply.MapC(4, func(_ context.Context, s string) (string, error) {
			time.Sleep(delay)
			return s + s, nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"testtest", "foofoo", "bazzbazz", "traintrain"})

		require.WithinDuration(t, start, time.Now(), delay+time.Millisecond*10)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.MapC(2, func(_ context.Context, _ string) (string, error) {
			return "", errTest
		}))

		testhelpers.RequireErroredPiping(t, out, errTest)
	})
}
