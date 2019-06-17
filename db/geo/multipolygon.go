package geo

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// for SRID:4326
type MultiPolygon []Polygon

type Polygon []LinearRing

type LinearRing []Point

func (mp *MultiPolygon) Scan(val interface{}) error {
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

	if err = CheckWKBType(r, byteOrder, WKBTypeMultiPolygon, WKBPointDimension, WKBSRID4326); err != nil {
		return err
	}

	var polygonNum uint32
	if err = binary.Read(r, byteOrder, &polygonNum); err != nil {
		return fmt.Errorf("Invalid PolygonNum(%d) in MultiPolygon", polygonNum)
	}
	*mp = make(MultiPolygon, int(polygonNum))
	for i := 0; i < int(polygonNum); i++ {
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
			return fmt.Errorf("Unsupported byte order %d in Polygons", wkbByteOrder)
		}

		var wkbGeometryType uint32
		if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
			return err
		}
		if wkbGeometryType&0xFF != WKBTypePolygon {
			return fmt.Errorf("Unsupported WKBType %d in MultiPolygons", wkbGeometryType)
		}

		var linearRingNum uint32
		if err = binary.Read(r, byteOrder, &linearRingNum); err != nil {
			return fmt.Errorf("Invalid LinearRingNum(%d) in Polygon", linearRingNum)
		}

		(*mp)[i] = make(Polygon, int(linearRingNum))
		for j := 0; j < int(linearRingNum); j++ {
			var pointNum uint32
			if err = binary.Read(r, byteOrder, &pointNum); err != nil {
				return fmt.Errorf("Invalid PointNum(%d) in LinearRing", pointNum)
			}

			(*mp)[i][j] = make(LinearRing, int(pointNum))

			for k := 0; k < int(pointNum); k++ {
				if err = binary.Read(r, byteOrder, &(*mp)[i][j][k]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (mp MultiPolygon) Value() (driver.Value, error) {
	if len(mp) == 0 {
		return "SRID=4326;MULTIPOLYGON EMPTY", nil
	}

	return "SRID=4326;MULTIPOLYGON" + mp.String(false), nil
}

func (mp MultiPolygon) String(hasPrefix bool) string {
	buf := bytes.NewBuffer(nil)

	for i, n := 0, len(mp); i < n; i++ {
		if i == 0 {
			buf.WriteString("(")
		}

		buf.WriteString(mp[i].String(hasPrefix))

		if i != n-1 {
			buf.WriteString(", ")
		} else {
			buf.WriteString(")")
		}
	}

	return buf.String()
}

func (p Polygon) String(hasPrefix bool) string {
	buf := bytes.NewBuffer(nil)

	for i, n := 0, len(p); i < n; i++ {
		if i == 0 {
			buf.WriteString("(")
		}

		buf.WriteString(p[i].String(hasPrefix))

		if i != n-1 {
			buf.WriteString(", ")
		} else {
			buf.WriteString(")")
		}
	}

	return buf.String()
}

func (l LinearRing) String(hasPrefix bool) string {
	buf := bytes.NewBuffer(nil)

	for i, n := 0, len(l); i < n; i++ {
		if i == 0 {
			buf.WriteString("(")
		}

		buf.WriteString(fmt.Sprintf("%v %v", l[i][0], l[i][1]))

		if i != n-1 {
			buf.WriteString(", ")
		} else {
			buf.WriteString(")")
		}
	}

	return buf.String()
}
