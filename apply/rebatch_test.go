package apply_test

import (
	"testing"

	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

func TestRebatch(t *testing.T) {
	t.Parallel()

	out := testhelpers.TestStageWith(t, apply.Rebatch(2), []any{[]string{"test", "foo", "bazz"}, []string{"train"}})

	testhelpers.RequireSuccessfulPiping(t, out, []any{[]any{"test", "foo"}, []any{"bazz", "train"}})
}
