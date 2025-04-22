package apply_test

import (
	"context"
	"testing"

	"github.com/zhulik/pips"
	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestZip(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.Zip(func(_ context.Context, s string) (string, error) {
			return s + s, nil
		}))

		// Expected pairs of original items and transformed results
		expected := []any{
			pips.NewP("test", "testtest"),
			pips.NewP("foo", "foofoo"),
			pips.NewP("bazz", "bazzbazz"),
			pips.NewP("train", "traintrain"),
		}

		testhelpers.RequireSuccessfulPiping(t, pips.OutChan[any](out), expected)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.Zip(func(_ context.Context, _ string) (string, error) {
			return "", errTest
		}))

		testhelpers.RequireErroredPiping(t, pips.OutChan[any](out), errTest)
	})
}
