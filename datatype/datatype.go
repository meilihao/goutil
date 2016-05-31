package datatype

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Strings []string

var (
	ZeroStrings               = Strings{}
	_           sql.Scanner   = &ZeroStrings
	_           driver.Valuer = ZeroStrings
)

func (ss Strings) Value() (driver.Value, error) {
	return ss.String(), nil
}

func (ss Strings) String() string {
	if len(ss) == 0 {
		return ""
	}
	return strings.Join(ss, ",")
}

func (ss *Strings) Scan(value interface{}) error {
	str := ""
	switch t := value.(type) {
	case string:
		str = t
	case []byte:
		str = string(t)
	default:
		typ := reflect.TypeOf(value)
		return errors.New(fmt.Sprintf("Strings.Scan cannot decode Type(%v)", typ.Kind()))
	}
	*ss = parseStrings(str)
	return nil
}

func parseStrings(str string) Strings {
	if str == "" {
		return Strings{}
	}
	return Strings(strings.Split(str, ","))
}
