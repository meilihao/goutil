package geo

import (
"testing"

"github.com/stretchr/testify/assert"
)

func TestIsValidPoint(t *testing.T) {
	point := Point{
		Lat: 38.804000,
		Lng: 121.353282,
	}
	assert.Equal(t, true, IsValidPointInChina(point))
}

