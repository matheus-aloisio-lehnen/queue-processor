package functions_test

import (
	"queue/core/infra/utils/functions"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIds(t *testing.T) {
	tests := map[string]struct {
		input        string
		expected     []int
		expectingErr bool
	}{
		"Lista válida com espaços": {
			input:        "1, 2, 3",
			expected:     []int{1, 2, 3},
			expectingErr: false,
		},
		"Lista válida sem espaços": {
			input:        "4,5,6",
			expected:     []int{4, 5, 6},
			expectingErr: false,
		},
		"Lista com valor não numérico": {
			input:        "1,abc,3",
			expected:     nil,
			expectingErr: true,
		},
		"String vazia": {
			input:        "",
			expected:     []int{},
			expectingErr: false,
		},
		"Lista com espaços extras": {
			input:        " 10 ,  20 ,30 ",
			expected:     []int{10, 20, 30},
			expectingErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := functions.ParseIds(tc.input)

			if tc.expectingErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
