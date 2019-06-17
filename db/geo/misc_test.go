package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPoint(t *testing.T) {
	point := Point{121.353282, 38.804000}
	assert.Equal(t, true, IsValidPointInChina(point))
}
