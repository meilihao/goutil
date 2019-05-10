// struct from [矢量空间数据格式](https://www.cnblogs.com/marsprj/archive/2013/02/08/2909452.html)
package geo

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	WKBTypePoint        = 0x01
	WKBTypePolygon      = 0x03
	WKBTypeMultiPolygon = 0x06
)

const (
	WKBPointDimension = 0x2000 // (lng,lat)
)

const (
	WKBSRID4326 = 0x000010E6
)

func CheckWKBType(r io.Reader, byteOrder binary.ByteOrder, wantGeometryType uint16, wantPointDimension uint16, wantGeometrySRID uint32) error {
	var wkbGeometryType uint16
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}
	if wkbGeometryType != wantGeometryType {
		return fmt.Errorf("Unsupported WKBType %d", wkbGeometryType&0xFF)
	}

	var wkbGeometryPointDimension uint16
	if err := binary.Read(r, byteOrder, &wkbGeometryPointDimension); err != nil {
		return err
	}
	if wkbGeometryPointDimension != wantPointDimension {
		return fmt.Errorf("Unsupported WKBPointDimension %d", wkbGeometryPointDimension)
	}

	var wkbGeometrySRID uint32
	if err := binary.Read(r, byteOrder, &wkbGeometrySRID); err != nil {
		return err
	}
	if wkbGeometrySRID != wantGeometrySRID {
		return fmt.Errorf("Unsupported SRID %d", wkbGeometrySRID)
	}

	return nil
}
