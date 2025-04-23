package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
		val, err := toStringSlice(value)
		if err != nil {
			panic(err.Error())
		}
		return val
	case []int:
		val, err := toIntegerSlice(value)
		if err != nil {
			panic(err.Error())
		}
		return val
	case map[string]any:
		val, err := toJsonMap(value)
		if err != nil {
			panic(err)
		}
		return val
	default:
		panic(fmt.Sprintf("Unsupported type: %v", reflect.TypeOf(def)))
	}
}

func toStringSlice(val any) ([]string, error) {
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("invalid data type not string")
	}
	commaValues := strings.Split(s, ",")
	if len(commaValues) < 1 {
		return nil, fmt.Errorf("comma missing in defining value")
	}
	var strArrayVal []string
	strArrayVal = append(strArrayVal, commaValues...)
	return strArrayVal, nil
}

func toIntegerSlice(val any) ([]int, error) {
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("invalid data type not integer")
	}
	commaValues := strings.Split(s, ",")
	if len(commaValues) < 1 {
		return nil, fmt.Errorf("comma missing in defining value")
	}
	var numIntArr []int

	for _, v := range commaValues {
		val, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		numIntArr = append(numIntArr, val)
	}
	return numIntArr, nil
}

func toJsonMap(value string) (map[string]any, error) {
	var jsonVal map[string]any
	err := json.Unmarshal([]byte(value), jsonVal)
	if err != nil {
		return nil, fmt.Errorf("unable to convert json into map:%v", err)
	}
	return jsonVal, nil
}
