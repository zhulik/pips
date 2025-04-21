package apply_test

import (
	"context"
	"testing"

	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestEach(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.Each(func(_ context.Context, _ string) error {
			return nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"test", "foo", "bazz", "train"})
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.Each(func(_ context.Context, _ string) error {
			return errTest
		}))

		testhelpers.RequireErroredPiping(t, out, errTest)
	})
}
