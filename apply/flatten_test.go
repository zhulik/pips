package apply_test

import (
	"testing"

	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestFlatten(t *testing.T) {
	t.Parallel()

	out := testhelpers.TestStageWith(t, apply.Flatten(), []any{[]any{"test", "foo", "bazz", "train"}})

	testhelpers.RequireSuccessfulPiping(t, out, []any{"test", "foo", "bazz", "train"})
}
