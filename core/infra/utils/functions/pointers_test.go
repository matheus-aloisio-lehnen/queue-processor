package functions_test

import (
	"github.com/stretchr/testify/assert"
	"queue/core/infra/utils/functions"
	"testing"
)

func TestPtrStr(t *testing.T) {
	str := "teste"
	ptr := functions.PtrStr(str)

	assert.NotNil(t, ptr)
	assert.Equal(t, str, *ptr)
}

func TestPtrBool(t *testing.T) {
	tests := map[string]bool{
		"Valor true":  true,
		"Valor false": false,
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := functions.PtrBool(val)

			assert.NotNil(t, ptr)
			assert.Equal(t, val, *ptr)
		})
	}
}
