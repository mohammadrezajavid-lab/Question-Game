package normalize_test

import (
	"github.com/stretchr/testify/assert"
	"golang.project/go-fundamentals/gameapp/pkg/normalize"
	"testing"
)

func TestNormalize_NormalizePhoneNumber(t *testing.T) {
	type addTest struct {
		arg      string
		expected string
	}

	var addTests = []addTest{
		{
			arg:      "+989196551929 ",
			expected: "+989196551929",
		},
		{
			arg:      " 09196551929",
			expected: "+989196551929",
		},
		{
			arg:      "9196551929",
			expected: "+989196551929",
		},
		{
			arg:      "989196551929",
			expected: "+989196551929",
		},
		{
			arg:      "89196551929",
			expected: "",
		},
		{
			arg:      "6551929+98919",
			expected: "",
		},
	}
	var strError = "invalid phone number format"

	for _, test := range addTests {
		res, err := normalize.New().NormalizePhoneNumber(test.arg)
		if err != nil {
			assert.EqualError(t, err, strError)
		}
		assert.Equal(t, test.expected, res)
	}
}
