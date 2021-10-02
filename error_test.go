package publicerror_test

import (
	"errors"
	"fmt"
	"testing"

	publicerror "github.com/xeoncross/public-error-go"
)

func TestFind(t *testing.T) {

	e := publicerror.Wrap(errors.New("danger"), "app problem", 0)

	testCases := []struct {
		desc   string
		err    error
		result error
	}{
		{
			desc:   "nil error",
			err:    nil,
			result: nil,
		},
		{
			desc:   "error with no publicerror",
			err:    errors.New("problem here"),
			result: nil,
		},
		{
			desc:   "error chain with no publicerror",
			err:    fmt.Errorf("level 2: %w", errors.New("problem here")),
			result: nil,
		},
		{
			desc:   "error chain with publicerror",
			err:    fmt.Errorf("level 2: %w", e),
			result: e,
		},
		{
			desc:   "error chain with publicerror used twice",
			err:    publicerror.Wrap(fmt.Errorf("level 2: %w", e), "say this instead", 0),
			result: publicerror.Error{Message: "say this instead"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := publicerror.Find(tc.err)
			if (tc.result == nil && got != nil) ||
				(tc.result != nil && tc.result.(publicerror.Error).Message != got.Message) {
				t.Errorf("unexpected result:\n\tgot: %v\n\twant: %v\n", got, tc.result)
			}
		})
	}
}
