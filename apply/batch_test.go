package apply_test

import (
	"testing"

	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestBatch(t *testing.T) {
	t.Parallel()

	out := testhelpers.TestStageWith(t, apply.Batch(2), []any{"test", "foo", "bazz", "train"})

	testhelpers.RequireSuccessfulPiping(t, out, []any{[]any{"test", "foo"}, []any{"bazz", "train"}})
}
