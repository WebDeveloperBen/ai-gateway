package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name           string
		input          ListRequest
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "Zero limit defaults to 100",
			input:          ListRequest{Limit: 0, Offset: 0},
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "Negative limit defaults to 100",
			input:          ListRequest{Limit: -1, Offset: 0},
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "Negative offset defaults to 0",
			input:          ListRequest{Limit: 50, Offset: -1},
			expectedLimit:  50,
			expectedOffset: 0,
		},
		{
			name:           "Valid values unchanged",
			input:          ListRequest{Limit: 25, Offset: 10},
			expectedLimit:  25,
			expectedOffset: 10,
		},
		{
			name:           "Both invalid values normalized",
			input:          ListRequest{Limit: -5, Offset: -10},
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "Large limit preserved",
			input:          ListRequest{Limit: 1000, Offset: 500},
			expectedLimit:  1000,
			expectedOffset: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePagination(tt.input)
			assert.Equal(t, tt.expectedLimit, result.Limit)
			assert.Equal(t, tt.expectedOffset, result.Offset)
		})
	}
}

func TestRepositoryBackend(t *testing.T) {
	t.Run("Constants defined correctly", func(t *testing.T) {
		assert.Equal(t, RepositoryBackend("postgres"), RepositoryPostgres)
		assert.Equal(t, RepositoryBackend("memory"), RepositoryMemory)
	})
}

func TestListRequest(t *testing.T) {
	t.Run("Struct fields", func(t *testing.T) {
		req := ListRequest{
			Limit:  50,
			Offset: 100,
		}
		assert.Equal(t, 50, req.Limit)
		assert.Equal(t, 100, req.Offset)
	})
}
