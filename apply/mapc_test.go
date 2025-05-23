package apply_test

import (
	"context"
	"strings"
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

		out := testhelpers.TestStringStage(t, apply.MapC(2, func(_ context.Context, s string) (string, error) {
			return s + s, nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"testtest", "foofoo", "bazzbazz", "traintrain"})
	})

	t.Run("runs concurrently", func(t *testing.T) {
		t.Parallel()

		start := time.Now()

		delay := 200 * time.Millisecond

		out := testhelpers.TestStringStage(t, apply.MapC(4, func(_ context.Context, s string) (string, error) {
			time.Sleep(delay)
			return s + s, nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"testtest", "foofoo", "bazzbazz", "traintrain"})

		require.WithinDuration(t, start, time.Now(), delay+time.Millisecond*10)
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestArrayStage(t, apply.MapC(2, func(_ context.Context, s []string) (string, error) {
			return strings.Join(s, " "), nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"a a", "b b", "c c", "d d"})
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStringStage(t, apply.MapC(2, func(_ context.Context, _ string) (string, error) {
			return "", errTest
		}))

		testhelpers.RequireErroredPiping(t, out, errTest)
	})

	t.Run("panic", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStringStage(t, apply.MapC(2, func(_ context.Context, _ string) (string, error) {
			panic(errTest)
		}))

		testhelpers.RequireErroredPiping(t, out, errTest)
	})
}
