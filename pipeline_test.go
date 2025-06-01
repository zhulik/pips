package pips_test

import (
	"context"
	"testing"

	"github.com/zhulik/pips"
	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

var (
	subPipe = apply.Pipeline(pips.New[string, string](apply.Map(func(_ context.Context, s string) (string, error) {
		return s + s, nil
	})))

	lenMap = apply.Map(func(_ context.Context, s string) (int, error) {
		return len(s), nil
	})

	nothingMap = apply.Map(func(_ context.Context, s []int) ([]int, error) {
		return s, nil
	})

	gt6Filter = apply.Filter(func(_ context.Context, i int) (bool, error) {
		return i > 6, nil
	})
)

func TestPipeline_Result(t *testing.T) {
	t.Parallel()

	out := pips.New[string, int]().
		Then(subPipe).
		Then(lenMap).
		Then(apply.Batch(3)).
		Then(nothingMap).
		Then(apply.Flatten()).
		Then(gt6Filter).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireSuccessfulPiping(t, out, []int{8, 8, 10})
}

func TestPipeline_NoResult(t *testing.T) {
	t.Parallel()

	out := pips.New[string, any]().
		Then(apply.Map(func(_ context.Context, _ string) (any, error) {
			return nil, nil
		})).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireSuccessfulPiping(t, out, []any{nil, nil, nil, nil})
}

func TestPipeline_StageError(t *testing.T) {
	t.Parallel()

	out := pips.New[string, any]().
		Then(apply.Map(func(_ context.Context, _ string) (any, error) {
			return nil, errTest
		})).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireErroredPiping(t, out, errTest)
}

func TestPipeline_StagePanic(t *testing.T) {
	t.Parallel()

	out := pips.New[string, any]().
		Then(apply.Map(func(_ context.Context, _ string) (any, error) {
			panic(errTest)
		})).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireErroredPiping(t, out, errTest)
}

func TestPipeline_StageOutputTypeCastingError(t *testing.T) {
	t.Parallel()

	out := pips.New[string, string]().
		Then(apply.Map(func(_ context.Context, _ string) (int, error) {
			return 0, nil
		})).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireErroredPiping(t, out, pips.ErrWrongType)
}

func TestPipeline_StageInputTypeCastingError(t *testing.T) {
	t.Parallel()

	out := pips.New[string, string]().
		Then(apply.Map(func(_ context.Context, _ int) (string, error) {
			return "", nil
		})).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireErroredPiping(t, out, pips.ErrWrongType)
}
