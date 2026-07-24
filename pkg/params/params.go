package params

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/IbrahimHYoussef/project-management/shared-utils-go/validator"
)

type QueryParams interface {
	validator.Validator
	MapQuery(*http.Request) (QueryParams, error)
}


func BindQuery[T any](r *http.Request) (*T, error) {
	// get the query dict
	values := r.URL.Query()

	// get distination type
	// var req T
	req := new(T)
	v := reflect.ValueOf(req)
	// make sure it is none nil pointer
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return nil, errors.New("dest Must be a None Nil Pointer")
	}
	// get element
	v = v.Elem()
	// get type
	t := v.Type()

	for i := range v.NumField() {
		// get element and type for the currnet itteration
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get json tag
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			continue
		}
		// ge the value of the tag
		val := values.Get(tag)
		if val == "" {
			continue
		}

		// set based on type
		switch field.Kind() {
		case reflect.String:
			field.SetString(val)
		case reflect.Int, reflect.Int64, reflect.Int32:
			num, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for '%s'", tag)
			}
			field.SetInt(int64(num))
		case reflect.Float32, reflect.Float64:
			num, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value for '%s'", tag)
			}
			field.SetFloat(float64(num))
		case reflect.Bool:
			bool, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid value for '%s'", tag)
			}
			field.SetBool(bool)
		case reflect.Slice:
			return nil,fmt.Errorf("unsupported field type: slices")
		default:
			return nil, fmt.Errorf("unsupported field type: %s", field.Kind())
		}
	}
	return req, nil
}
