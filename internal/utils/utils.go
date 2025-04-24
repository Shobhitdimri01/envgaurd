package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ReplaceEnvPlaceholders(pattern *regexp.Regexp, input string) string {
	matches := pattern.FindAllStringSubmatchIndex(input, -1)
	begin := 0
	var result strings.Builder
	for _, match := range matches {
		startIdx := match[0] //matches full index
		endIdx := match[1]   // by grouping - ${}
		result.WriteString(input[begin:startIdx])
		key := strings.TrimSpace(input[match[2]:match[3]]) // matches key only
		envVal, ok := os.LookupEnv(key)
		if !ok {
			panic(fmt.Sprintf("unable to find env variable key named %s", input[startIdx:endIdx]))
		}
		result.WriteString(envVal)
		begin = endIdx
	}
	result.WriteString(input[begin:])
	return result.String()
}
