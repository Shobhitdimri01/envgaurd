package validate

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func parsetoStringSlice(val any) ([]string, error) {
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

func parsetoIntegerSlice(val any) ([]int, error) {
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

func parsetoJsonMap(value string) (map[string]any, error) {
	var jsonVal map[string]any
	err := json.Unmarshal([]byte(value), &jsonVal)
	if err != nil {
		return nil, fmt.Errorf("unable to convert json into map:%v", err)
	}
	return jsonVal, nil
}