package validate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Shobhitdimri01/envgaurd/internal/utils"
)

func ParseEnvFile(path string, set func(key, val string)) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue //invalid format, skip
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		//check for inline comment
		val = utils.StripInlineComment(val)
		set(key, val)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}
	return nil
}

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
