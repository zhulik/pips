package pips_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhulik/pips"
)

var (
	errTest = errors.New("test error")
)

func readAll[T any](ch <-chan T) []T {
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

func inputChan() <-chan pips.D[string] {
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

func testPipeline(t *testing.T, stages ...pips.Stage) <-chan pips.D[string] {
	t.Helper()

	return pips.New[string, string](stages...).Run(t.Context(), inputChan())
}

func requireSuccessfulPiping[T any](t *testing.T, out <-chan pips.D[T], expected []T) {
	collected := readAll(out)

	require.Equal(t, expected,
		Map(collected, func(item pips.D[T], _ int) T {
			require.NoError(t, item.Error())
			return item.Value()
		}),
	)
}

func requireErroredPiping[T any](t *testing.T, out <-chan pips.D[T], err error) {
	collected := readAll(out)

	require.Len(t, collected, 1)
	require.ErrorIs(t, collected[0].Error(), err)
}
