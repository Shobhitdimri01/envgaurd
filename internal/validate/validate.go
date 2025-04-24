package validate

import (
	"fmt"
	"reflect"
	"strconv"
)

func Validate(value string, def any) any {
	switch def.(type) {
	case int:
		val, err := strconv.Atoi(value)
		if err != nil {
			panic("invalid value expected Integer type ")
		}
		return val
	case bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			panic("invalid value expected Boolean type ")
		}
		return val
	case string:
		return value
	case float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic("invalid value expected float64 type")
		}
		return val
	case []string:
		val, err := parsetoStringSlice(value)
		if err != nil {
			panic(err.Error())
		}
		return val
	case []int:
		val, err := parsetoIntegerSlice(value)
		if err != nil {
			panic(err.Error())
		}
		return val
	case map[string]any:
		val, err := parsetoJsonMap(value)
		if err != nil {
			panic(err)
		}
		return val
	default:
		panic(fmt.Sprintf("Unsupported type: %v", reflect.TypeOf(def)))
	}
}
