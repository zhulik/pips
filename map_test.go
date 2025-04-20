package pips_test

import (
	"context"
	"testing"

	"github.com/zhulik/pips/apply"
)

var (
	duplicateMap = apply.Map(func(_ context.Context, s string) (string, error) {
		return s + s, nil
	})

	erroredMap = apply.Map(func(_ context.Context, _ string) (string, error) {
		return "", errTest
	})
)

func TestPipeline_Map(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		out := testPipeline(t, duplicateMap)

		requireSuccessfulPiping(t, out, []string{"testtest", "foofoo", "bazzbazz", "traintrain"})
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		out := testPipeline(t, erroredMap)

		requireErroredPiping(t, out, errTest)
	})
}
