package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhulik/pips"
)

func InputChan() <-chan pips.D[string] {
	ch := make(chan pips.D[string])

	go func() {
		ch <- pips.NewD("test")
		ch <- pips.NewD("foo")
		ch <- pips.NewD("bazz")
		ch <- pips.NewD("train")
		close(ch)
	}()

	return ch
}

func TestStage(t *testing.T, stage pips.Stage) <-chan pips.D[any] {
	t.Helper()

	return TestStageWith(t, stage, []any{"test", "foo", "bazz", "train"})
}

func TestStageWith(t *testing.T, stage pips.Stage, items []any) <-chan pips.D[any] {
	t.Helper()

	ch := make(chan pips.D[any])

	go func() {
		for _, item := range items {
			ch <- pips.AnyD(item)
		}
		close(ch)
	}()

	out := make(chan pips.D[any])
	go func() {
		stage(t.Context(), ch, out)
		close(out)
	}()

	return out
}

func RequireSuccessfulPiping[T any](t *testing.T, out pips.OutChan[T], expected []T) {
	t.Helper()

	collected, err := out.Collect(t.Context())
	require.NoError(t, err)

	require.Equal(t, expected, collected)
}

func RequireErroredPiping[T any](t *testing.T, out pips.OutChan[T], expected error) {
	t.Helper()

	_, err := out.Collect(t.Context())

	require.ErrorIs(t, err, expected)
}
