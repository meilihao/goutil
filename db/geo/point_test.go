package geo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointStringRepresentation(t *testing.T) {
	point := Point{-121.353282, 38.804000}
	assert.Equal(t, "SRID=4326;POINT(-121.353282 38.804)", point.String())
}

func TestPointValue(t *testing.T) {
	point := Point{-121.353282, 38.804000}
	v, err := point.Value()
	assert.Nil(t, err)
	assert.Equal(t, "SRID=4326;POINT(-121.353282 38.804)", v)
}

func TestPointScanning(t *testing.T) {
	raw := "0101000020E6100000E6CE4C309C565EC023827170E9664340"
	point := &Point{}
	if err := point.Scan([]byte(raw)); err != nil {
		assert.Error(t, err)
	}

	fmt.Println(*point)
	assert.Equal(t, 38.803999, (*point)[1])
	assert.Equal(t, -121.353283, (*point)[0])
}

func TestScanNilPoint(t *testing.T) {
	var np NullPoint
	err := np.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, np.Valid)
}

func TestNullPointScanning(t *testing.T) {
	raw := "0101000020E6100000E6CE4C309C565EC023827170E9664340"
	np := &NullPoint{}
	if err := np.Scan([]byte(raw)); err != nil {
		assert.Error(t, err)
	}

	assert.True(t, np.Valid)
	assert.Equal(t, 38.803999, np.Point[1])
	assert.Equal(t, -121.353283, np.Point[0])
}
