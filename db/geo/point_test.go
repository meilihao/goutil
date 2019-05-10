package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointStringRepresentation(t *testing.T) {
	point := Point{
		Lat: 38.804000,
		Lng: -121.353282,
	}
	assert.Equal(t, "SRID=4326;POINT(-121.353282 38.804)", point.String())
}

func TestPointValue(t *testing.T) {
	point := Point{
		Lat: 38.804000,
		Lng: -121.353282,
	}
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

	assert.Equal(t, 38.803999, point.Lat)
	assert.Equal(t, -121.353283, point.Lng)
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
	assert.Equal(t, 38.803999, np.Point.Lat)
	assert.Equal(t, -121.353283, np.Point.Lng)
}