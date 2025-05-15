package apply_test

import (
	"context"
	"testing"

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

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStage(t, apply.MapC(2, func(_ context.Context, _ string) (string, error) {
			return "", errTest
		}))

		testhelpers.RequireErroredPiping(t, out, errTest)
	})
}
