package apply_test

import (
	"context"
	"strings"
	"testing"

	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestMap(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStringStage(t, apply.Map(func(_ context.Context, s string) (string, error) {
			return s + s, nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"testtest", "foofoo", "bazzbazz", "traintrain"})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestArrayStage(t, apply.Map(func(_ context.Context, s []string) (string, error) {
			return strings.Join(s, " "), nil
		}))

		testhelpers.RequireSuccessfulPiping(t, out, []any{"a a", "b b", "c c", "d d"})
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testhelpers.TestStringStage(t, apply.Map(func(_ context.Context, _ string) (string, error) {
			return "", errTest
		}))

		testhelpers.RequireErroredPiping(t, out, errTest)
	})
}
