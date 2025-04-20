package apply_test

import (
	"context"
	"testing"

	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.Filter(func(_ context.Context, s string) (bool, error) {
			return len(s) > 3, nil
		}))
		testhelpers.RequireSuccessfulPiping(t, out, []any{"test", "bazz", "train"})
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.Filter(func(_ context.Context, _ string) (bool, error) {
			return true, errTest
		}))

		testhelpers.RequireErroredPiping(t, out, errTest)
	})
}
