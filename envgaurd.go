package envgaurd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Shobhitdimri01/envgaurd/internal/utils"
	"github.com/Shobhitdimri01/envgaurd/internal/validate"
)

var sensitiveInfo = &sensitiveData{
	maskingKeys: make(map[string]*envMetaData),
}

type sensitiveData struct {
	maskingKeys map[string]*envMetaData
}

type envMetaData struct {
	value any
	mask  bool
}

func setEnv(tempEnv map[string]string) {
	for key, val := range tempEnv {
		os.Setenv(key, val)
		sensitiveInfo.maskingKeys[strings.ToUpper(key)] = &envMetaData{value: val}
	}
}

// Load manually loads the .env file into environment variables
// Only sets environment variable if not already set it doesn't override
func Load(path string) error {
	tempEnv := make(map[string]string)
	err := validate.ParseEnvFile(path, func(key, val string) {
		value := os.Getenv(key)
		if value == "" {
			tempEnv[key] = val
		} else {
			tempEnv[key] = value
		}
	})
	if err != nil {
		return err
	}
	setEnv(tempEnv)
	return nil
}

// OverLoad overwrite all the existing environment variable
func OverLoad(path string) error {
	tempEnv := make(map[string]string)
	err := validate.ParseEnvFile(path, func(key, val string) {
		tempEnv[key] = val
	})
	if err != nil {
		return err
	}
	sensitiveInfo.clear()
	setEnv(tempEnv)
	return nil
}

// LoadFromFileWithValidation loads environment variables from a file and validates them
func LoadFromFileWithValidation(path string, requiredKeys []string) error {
	tempEnv := make(map[string]string)

	// Step 1: Parse the file into a temporary map
	err := validate.ParseEnvFile(path, func(key, val string) {
		tempEnv[key] = val
	})
	if err != nil {
		panic(fmt.Sprintf("Error reading env file: %v", err))
	}
	// Step 2: Validate required keys
	for _, key := range requiredKeys {
		if val, ok := tempEnv[key]; !ok {
			panic(fmt.Sprintf("Missing required environment variable: %s", key))
		} else if val == "" {
			panic(fmt.Sprintf("Missing required value for key: %s", key))
		}
	}
	// Step 3: Set the environment variables
	setEnv(tempEnv)
	return nil
}

// Require checks that an env variable is set, and panics if not.
// The Require function is designed to enforce that certain environment variables are always present when your application runs.
func Required(key string) {
	if os.Getenv(key) == "" {
		panic(fmt.Sprintf("Missing required env var: %s", key))
	}
}

//GetInt extract Value from env dot file and return Integer Value if not given returns default value
func GetInt(key string, defaultVal int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	val := validate.Validate(value, defaultVal)
	IntVal, ok := val.(int)
	if !ok {
		return defaultVal // fallback if type assertion fails
	}
	return IntVal
}

// GetStr extracts a value from the environment variables and returns it as a string.
// If the key is not found or the value is empty, it returns the provided default.
func GetStr(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	val := validate.Validate(value, defaultVal)
	strVal, ok := val.(string)
	if !ok {
		return defaultVal // fallback if type assertion fails
	}
	return strVal
}

func GetBool(key string, defaultVal bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	val := validate.Validate(value, defaultVal)
	BoolVal, ok := val.(bool)
	if !ok {
		return defaultVal // fallback if type assertion fails
	}
	return BoolVal
}

// GetFloat retrieves a float value from the environment, with a default if not found
func GetFloat64(key string, defaultVal float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	val := validate.Validate(value, defaultVal)
	FloatVal, ok := val.(float64)
	if !ok {
		return defaultVal // fallback if type assertion fails
	}
	return FloatVal
}

// GetStringArray retrieves a comma-separated list of strings from the environment
// and returns them as a slice of strings
func GetStringArray(key string, defaultVal []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	val := validate.Validate(value, defaultVal)
	strArrayVal, ok := val.([]string)
	if !ok {
		return defaultVal // fallback if type assertion fails
	}
	return strArrayVal
}

// GetIntegerArray retrieves a comma-separated list of strings from the environment
// and returns them as a slice of integer
func GetIntegerArray(key string, defaultVal []int) []int {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	val := validate.Validate(value, defaultVal)
	intArrayVal, ok := val.([]int)
	if !ok {
		return defaultVal // fallback if type assertion fails
	}
	return intArrayVal
}

// GetEnvAsMap retrieves a JSON-encoded map from an env var, or returns default if not found or invalid
func GetEnvAsMap(key string, defaultVal map[string]any) map[string]any {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	val := validate.Validate(key, defaultVal)
	mapValues, ok := val.(map[string]any)
	if !ok {
		return defaultVal // fallback if type assertion fails
	}
	return mapValues
}

// ReplaceEnvPlaceholders replaces ${VAR_NAME} in the input string with its corresponding env value
func GetPlaceHolderValue(key string, defaultVal string) string {
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)\}`)
	val := os.Getenv(key)
	// If no placeholder pattern is found, panic (since that’s the expected behavior here)
	if !re.MatchString(val) {
		panic(fmt.Sprintf("value for key '%s' doesn't match placeholder syntax ${some-value}: %s", key, val))
	}
	// if exist now replace n place holders
	result := utils.ReplaceEnvPlaceholders(re, val)
	return result
}

// PrintEnvVars is a utility function that iterates through all the environment variables of the current process and prints them in a formatted manner. For any environment variable deemed sensitive (e.g., containing secrets, tokens, passwords),
// its value is masked before printing to prevent leaking sensitive information.
func PrintEnvVars() {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := strings.TrimSpace(parts[0])
		envMetaData, ok := sensitiveInfo.maskingKeys[strings.ToUpper(key)]
		if ok {
			val := envMetaData.value
			if envMetaData.mask {
				val = maskValue(envMetaData.value)
			}
			fmt.Printf("%-2s = %s\n", key, val)
		}
	}
}

// Masking lets users define which key or keys are considered sensitive will be shown encrypted(eg.xxxx) in console
func Masking(keys ...string) {
	if sensitiveInfo.maskingKeys == nil {
		sensitiveInfo.maskingKeys = make(map[string]*envMetaData)
	}
	for _, v := range keys {
		k := strings.TrimSpace(v)
		val, ok := os.LookupEnv(k)
		if !ok {
			panic(fmt.Sprintf("key with name:%s not found in .env file\n", k))
		}
		sensitiveInfo.maskingKeys[strings.ToUpper(k)] = &envMetaData{
			value: val,
			mask:  true,
		}
	}
}

// Masks the value for printing
func maskValue(val any) string {
	switch val := val.(type) {
	case string:
		v := val
		if len(v) <= 4 {
			return strings.Repeat("*", len(v))
		}
		return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
	case int, float64, bool:
		return "***"
	default:
		return "******"
	}
}

func (s *sensitiveData) clear() {
	s.maskingKeys = make(map[string]*envMetaData)
}

// ResetEnv clears all environment variables for the current process.
// It also removes all tracked sensitive values.
// Returns an error if any call to os.Unsetenv fails.
func ResetEnv() error {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			if err := os.Unsetenv(parts[0]); err != nil {
				return err
			}
		}
	}
	//reseting map
	sensitiveInfo.maskingKeys = map[string]*envMetaData{}
	return nil
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
