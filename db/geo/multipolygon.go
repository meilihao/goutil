package geo

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// for SRID:4326
type MultiPolygon struct {
	Polygons []Polygon
}

type Polygon struct {
	LinearRings []LinearRing
}

type LinearRing struct {
	Points []Point
}

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
	mp.Polygons = make([]Polygon, int(polygonNum))
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

		mp.Polygons[i].LinearRings = make([]LinearRing, int(linearRingNum))
		for j := 0; j < int(linearRingNum); j++ {
			var pointNum uint32
			if err = binary.Read(r, byteOrder, &pointNum); err != nil {
				return fmt.Errorf("Invalid PointNum(%d) in LinearRing", pointNum)
			}

			mp.Polygons[i].LinearRings[j].Points = make([]Point, int(pointNum))

			for k := 0; k < int(pointNum); k++ {
				if err = binary.Read(r, byteOrder, &mp.Polygons[i].LinearRings[j].Points[k]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
