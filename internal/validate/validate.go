package validate

import (
	"fmt"
	"reflect"
	"strconv"
)

func Validate(value string, def interface{}) interface{} {
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
	default:
		panic(fmt.Sprintf("Unsupported type: %v", reflect.TypeOf(def)))
	}
}
