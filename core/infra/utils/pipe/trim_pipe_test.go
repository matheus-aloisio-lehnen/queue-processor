package pipe_test

import (
	"github.com/stretchr/testify/assert"
	"queue/core/infra/utils/pipe"
	"testing"
	"time"
)

type ComplexDTO struct {
	Name      string
	Email     string
	Age       int
	Active    bool
	Balance   float64
	Tags      []string
	Metadata  map[string]string
	CreatedAt time.Time
}

func TestTrimStrings_ShouldTrimOnlyStrings(t *testing.T) {
	dto := ComplexDTO{
		Name:      "   Matheus   ",
		Email:     "   matheus@example.com   ",
		Age:       30,
		Active:    true,
		Balance:   100.5,
		Tags:      []string{"  a  ", " b "},
		Metadata:  map[string]string{"key": "   value   "},
		CreatedAt: time.Now(),
	}
	pipe.TrimPipe(&dto)
	assert.Equal(t, "Matheus", dto.Name)
	assert.Equal(t, "matheus@example.com", dto.Email)
	assert.Equal(t, 30, dto.Age)
	assert.True(t, dto.Active)
	assert.Equal(t, 100.5, dto.Balance)
	assert.Equal(t, []string{"  a  ", " b "}, dto.Tags)
	assert.Equal(t, "   value   ", dto.Metadata["key"])
	assert.NotZero(t, dto.CreatedAt)
}
