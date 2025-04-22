package envgaurd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Shobhitdimri01/envgaurd/internal/validate"
)

// Load manually loads the .env file into environment variables
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
		os.Setenv(parts[0], parts[1])
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

func GetInt(key string, def int) int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	val := validate.Validate(value, def)
	return val.(int)
}
