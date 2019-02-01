package mapper

import (
	"errors"
	"reflect"
)

type ColScanner interface {
	Columns() ([]string, error)
	Scan(dest ...interface{}) error
	Err() error
}

func MapScan(r ColScanner, dest map[string]interface{}) error {
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
}

func StructScan(rows ColScanner, dest interface{}) (err error) {
	destValue := reflect.ValueOf(dest)
	elemType := destValue.Type()

	if elemType.Kind() != reflect.Ptr {
		return errors.New("slice elem must ptr ")
	}

	rowMap := make(map[string]interface{})
	if err = MapScan(rows, rowMap); err != nil {
		return
	}
	if err = MapperMap(rowMap, dest); err != nil {
		return
	}

	return
}
