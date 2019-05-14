// Based on https://github.com/jinzhu/gorm/issues/142
// Point NullPoint fork from https://github.com/dewski/spatial
package geo

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// for SRID:4326
type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

func (p Point) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Lng, p.Lat)
}

// [GIS基本概念](https://blog.csdn.net/alinshen/article/details/78503333)
func (p *Point) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err = binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("Unsupported byte order %d", wkbByteOrder)
	}

	if err = CheckWKBType(r, byteOrder, WKBTypePoint, WKBPointDimension, WKBSRID4326); err != nil {
		return err
	}

	if err = binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}

type NullPoint struct {
	Point Point
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullPoint) Scan(val interface{}) error {
	if val == nil {
		n.Point, n.Valid = Point{}, false
		return nil
	}

	n.Valid = true
	return n.Point.Scan(val)
}

// Value implements the driver Valuer interface.
func (n NullPoint) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Point, nil
}
