package envgaurd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

const filepath = "testdata/.env"

func TestLoad(t *testing.T) {
	//cleaning values
	defer ResetEnv()
	// Set TEST_KEY1 before loading to ensure it doesn't get overridden
	err := os.Setenv("APP_NAME", "demo-app")
	if err != nil {
		t.Fatalf("failed to set APP_NAME: %v", err)
	}
	err = os.Setenv("APP_ENV", "prod")
	if err != nil {
		t.Fatalf("failed to set APP_NAME: %v", err)
	}
	//load doesn't override the existing key
	err = Load(filepath)
	if err != nil {
		t.Fatalf("failed to load file: %v", err)
	}
	if val := os.Getenv("APP_NAME"); val != "demo-app" {
		t.Errorf("expected APP_NAME to be 'demo-app', got '%s'", val)
	}

	// Check TEST_KEY2 is set from file
	if val := os.Getenv("APP_ENV"); val != "prod" {
		t.Errorf("expected APP_ENV to be 'development', got '%s'", val)
	}
}
func TestOverLoad(t *testing.T) {
	// Set TEST_KEY1 before loading to ensure it doesn't get overridden
	err := os.Setenv("APP_NAME", "demo-app")
	if err != nil {
		t.Fatalf("failed to set APP_NAME: %v", err)
	}
	//value not set
	err = os.Unsetenv("APP_ENV")
	if err != nil {
		t.Fatalf("failed to unset APP_NAME: %v", err)
	}

	//cleaning values
	defer ResetEnv()

	//load doesn't override the existing key
	err = OverLoad(filepath)
	if err != nil {
		t.Fatalf("failed to load file: %v", err)
	}
	if val := os.Getenv("APP_NAME"); val != "My Cool App" {
		t.Errorf("expected APP_NAME to be 'demo-app', got '%s'", val)
	}

	// Check TEST_KEY2 is set from file
	if val := os.Getenv("APP_ENV"); val != "development" {
		t.Errorf("expected APP_ENV to be 'development', got '%s'", val)
	}
}

func TestSetEnv(t *testing.T) {
	// Setup: initialize maskingKeys map
	sensitiveInfo.maskingKeys = make(map[string]*envMetaData)

	tempEnv := map[string]string{
		"APP_NAME": "demo-app",
	}
	defer ResetEnv()

	setEnv(tempEnv)

	if val := os.Getenv("APP_NAME"); val != "demo-app" {
		t.Errorf("expected APP_NAME to be 'demo-app', got '%s'", val)
	}

	info, ok := sensitiveInfo.maskingKeys["APP_NAME"]
	if !ok || info == nil {
		t.Errorf("expected entry in maskingKeys for APP_NAME, found none")
	} else {
		if info.value != "demo-app" {
			t.Errorf("expected value to be 'demo-app', got '%s'", info.value)
		}
		if info.mask {
			t.Errorf("expected mask to be false, got true")
		}
	}
	sensitiveInfo.clear()
}

func TestLoadFromFileWithValidation(t *testing.T) {
	testcaseTable := []struct {
		requiredKeys []string
		panicMsg     string
	}{
		{
			[]string{"APP_NAME", "APP_PORT"},
			"",
		},
		{
			[]string{"APP_NAME", "APP_PORT", "Prod_Engine"},
			"Missing required environment variable: Prod_Engine",
		},
		{
			[]string{"APP_NAME", "APP_PORT", "DB_PATH"},
			"Missing required value for key: DB_PATH",
		},
	}

	for i, j := range testcaseTable {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					msg := fmt.Sprint(r)
					if msg != j.panicMsg {
						t.Errorf("got unexpected panic msg: %s, expected: %s", msg, j.panicMsg)
					}
				}
			}()

			defer ResetEnv()

			if err := LoadFromFileWithValidation(filepath, j.requiredKeys); err != nil {
				t.Errorf("got error with file: %v", err)
			}
		})
	}
}

func TestPrintEnvVars(t *testing.T) {
	// Set up mock environment variables
	defer ResetEnv()
	_ = os.Setenv("DB_PASSWORD", "supersecret")
	_ = os.Setenv("APP_ENV", "production")
	// Set up masking info
	sensitiveInfo.maskingKeys = map[string]*envMetaData{
		"DB_PASSWORD": {value: "supersecret", mask: true},
		"APP_ENV":     {value: "production", mask: false},
	}

	expected := []string{
		fmt.Sprintf("%-2s = %s", "DB_PASSWORD", "su*******et"), // 11 asterisks for "supersecret"
		fmt.Sprintf("%-2s = %s", "APP_ENV", "production"),
	}

	original := os.Stdout
	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintEnvVars()
	w.Close()
	os.Stdout = original
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	consoleOutput := buf.String()
	res := strings.Split(consoleOutput, "\n")
	for i := 0; i < 2; i++ {
		if res[i] != expected[i] {
			t.Errorf("console output expected %s but got %s", expected[i], res[i])
		}
	}
}
