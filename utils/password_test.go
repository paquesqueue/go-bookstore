//go:build unit
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	t.Run("TestPasswordShouldReturnNoError", func(t *testing.T) {
		// Arrange
		mockData := "123456"
		// Act
		result, err := HashPassword(mockData)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.NotEqual(t, mockData, result)
	})

	t.Run("TestPasswordShouldReturnError", func(t *testing.T) {
		// Arrange
		mockData := ""

		// Act
		result, err := HashPassword(mockData)

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, "", result)
	})
}
