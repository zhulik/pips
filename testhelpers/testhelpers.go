package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhulik/pips"
)

func ReadAll[T any](ch <-chan T) []T {
	collection := []T{}

	for item := range ch {
		collection = append(collection, item)
	}

	return collection
}

func Map[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i := range collection {
		result[i] = iteratee(collection[i], i)
	}

	return result
}

func InputChan() <-chan pips.D[any] {
	ch := make(chan pips.D[any])

	go func() {
		ch <- pips.AnyD("test")
		ch <- pips.AnyD("foo")
		ch <- pips.AnyD("bazz")
		ch <- pips.AnyD("train")
		close(ch)
	}()

	return ch
}

func TestStage(t *testing.T, stage pips.Stage) <-chan pips.D[any] {
	t.Helper()

	out := make(chan pips.D[any])
	go func() {
		stage.Run(t.Context(), InputChan(), out)
		close(out)
	}()

	return out
}

func RequireSuccessfulPiping[T any](t *testing.T, out <-chan pips.D[T], expected []T) {
	collected := ReadAll(out)

	require.Equal(t, expected,
		Map(collected, func(item pips.D[T], _ int) T {
			require.NoError(t, item.Error())
			return item.Value()
		}),
	)
}

func RequireErroredPiping[T any](t *testing.T, out <-chan pips.D[T], err error) {
	collected := ReadAll(out)

	require.Len(t, collected, 1)
	require.ErrorIs(t, collected[0].Error(), err)
}
