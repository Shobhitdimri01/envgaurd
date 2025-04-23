package envgaurd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Shobhitdimri01/envgaurd/internal/validate"
)

// Load manually loads the .env file into environment variables
// Only sets environment variable if not already set it doesn't override
func Load(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("unable to open .env file:%v", err)
	}
	defer file.Close()
	scanLine := bufio.NewScanner(file)
	for scanLine.Scan() {
		line := scanLine.Text()
		// Ignore comments or empty lines
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		//invalid key-value
		if len(parts) != 2 {
			continue //invalid format, skip
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if os.Getenv(key) == "" { // Only set if not already set
			os.Setenv(key, val)
		}

	}
	if err := scanLine.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}
	return nil
}

// Require checks that an env variable is set, and panics if not.
// The Require function is designed to enforce that certain environment variables are always present when your application runs.
func Required(key string) {
	if os.Getenv(key) == "" {
		panic(fmt.Sprintf("Missing required env var: %s", key))
	}
}

// OverLoad overwrite all the existing environment variable
func OverLoad(filepath string) error {
	file, err := os.Open(filepath)
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
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		os.Setenv(key, val)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}
	return nil
}

// GetInt extract Value from env dot file and return Integer Value if not given returns default value
func GetInt(key string, def int) int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(value, def)
	IntVal, ok := val.(int)
	if !ok {
		return def // fallback if type assertion fails
	}
	return IntVal
}

// GetStr extracts a value from the environment variables and returns it as a string.
// If the key is not found or the value is empty, it returns the provided default.
func GetStr(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(value, def)
	strVal, ok := val.(string)
	if !ok {
		return def // fallback if type assertion fails
	}
	return strVal
}

func GetBool(key string, def bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(value, def)
	BoolVal, ok := val.(bool)
	if !ok {
		return def // fallback if type assertion fails
	}
	return BoolVal
}

// GetFloat retrieves a float value from the environment, with a default if not found
func GetFloat64(key string, def float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(value, def)
	FloatVal, ok := val.(float64)
	if !ok {
		return def // fallback if type assertion fails
	}
	return FloatVal
}

// GetStringArray retrieves a comma-separated list of strings from the environment
// and returns them as a slice of strings
func GetStringArray(key string, def []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(value, def)
	strArrayVal, ok := val.([]string)
	if !ok {
		return def // fallback if type assertion fails
	}
	return strArrayVal
}

// GetIntegerArray retrieves a comma-separated list of strings from the environment
// and returns them as a slice of integer
func GetIntegerArray(key string, def []int) []int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(value, def)
	intArrayVal, ok := val.([]int)
	if !ok {
		return def // fallback if type assertion fails
	}
	return intArrayVal
}

// GetEnvAsMap retrieves a JSON-encoded map from an env var, or returns default if not found or invalid
func GetEnvAsMap(key string, def map[string]any) map[string]any {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(key, def)
	mapValues, ok := val.(map[string]any)
	if !ok {
		return def // fallback if type assertion fails
	}
	return mapValues
}

// LoadFromFileWithValidation loads environment variables from a file and validates them
func LoadFromFileWithValidation(filepath string, requiredKeys []string) error {
	if err := Load(filepath); err != nil {
		return err
	}
	for _, key := range requiredKeys {
		Required(key) // Check if required keys are present
	}
	return nil
}

// PrintEnvVars is a utility function that iterates through all the environment variables of the current process and prints them in a formatted manner. For any environment variable deemed sensitive (e.g., containing secrets, tokens, passwords),
// its value is masked before printing to prevent leaking sensitive information.
func PrintEnvVars() {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		val := parts[1]
		if isSensitive(key) {
			val = maskValue(val)
		}
		fmt.Printf("%-20s = %s\n", key, val)
	}
}

// Checks if the key should be masked in log output
func isSensitive(key string) bool {
	key = strings.ToUpper(key)
	for k := range MaskingKeys {
		if strings.Contains(key, k) {
			return true
		}
	}
	return false
}

var MaskingKeys = make(map[string]bool)

// Masking lets users define which keys are considered sensitive will be shown encrypted(eg.xxxx)
func Masking(keys ...string) {
	for _, key := range keys {
		MaskingKeys[strings.ToUpper(key)] = true
	}
}

// Masks the value for printing
func maskValue(val string) string {
	if len(val) <= 4 {
		return strings.Repeat("*", len(val))
	}
	return val[:2] + strings.Repeat("*", len(val)-4) + val[len(val)-2:]
}

/*
Problem | Default Go | envgaurd Vision
Type safety | ❌ | ✅ GetInt, GetBool
Validation | ❌ | ✅ Required/Optional checks
Defaults | ❌ | ✅ Inline defaults
.env loading | ❌ | ✅ Built-in
Security masking | ❌ | ✅ Sensitive var protection
Multiple environments | ❌ | ✅ Profiles/namespaces
Centralized error handling | ❌ | ✅ All at once validation

*/
