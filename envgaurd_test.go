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

func TestGetInt(t *testing.T) {
	testTable := []struct {
		tcName      string
		key         string
		value       string
		defaultVal  int
		expectedRes any
	}{
		{"tc01", "Port", "8080", 0, 8080},                                 //passing tc
		{"tc02", "Port", "abc", 0, "invalid value expected Integer type"}, // should panic
		{"tc03", "Port", "", 5600, 5600},                                  //should take default value
		{"tc04", "Port", "8080", 3441, 8080},                              // shouldn't take default value
	}

	for _, j := range testTable {
		t.Run(j.tcName, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					err := fmt.Sprint(r)
					if err != j.expectedRes {
						t.Errorf("expected res: %s but got %s", j.expectedRes, err)
					}
				}
			}()

			if err := os.Setenv(j.key, j.value); err != nil {
				t.Error(err)
			}

			defer os.Unsetenv(j.key)
			val := GetInt(j.key, j.defaultVal)
			if j.key == "" {
				if j.defaultVal != val {
					t.Errorf("Expected result %d but got %d", j.defaultVal, val)
				}
			} else {
				if val != j.expectedRes {
					t.Errorf("Expected result %d but got %d", j.expectedRes, val)
				}
			}

		})
	}
}
func TestGetStr(t *testing.T) {
	testTable := []struct {
		tcName      string
		key         string
		value       string
		defaultVal  string
		expectedRes any
	}{
		{"tc01", "module_name", "envgaurd", "", "envgaurd"},            //passing tc
		{"tc02", "module_name", "123", "", "123"},                      //should be string
		{"tc03", "module_name", "", "envgaurd-v2", "envgaurd-v2"},      //should take default value
		{"tc04", "module_name", "envgaurd", "envgaurd-v2", "envgaurd"}, // shouldn't take default value
	}

	for _, j := range testTable {
		if err := os.Setenv(j.key, j.value); err != nil {
			t.Error(err)
		}

		defer os.Unsetenv(j.key)
		val := GetStr(j.key, j.defaultVal)
		if j.key == "" {
			if j.defaultVal != val {
				if j.defaultVal == "" {
					t.Errorf(`tcname : %s Expected result {""} but got {%s}`, j.tcName, val)
				} else {
					t.Errorf(`tcname : %s Expected result {%s} but got {%s}`, j.tcName, j.defaultVal, val)
				}

			}
		} else {
			if val != j.expectedRes {
				if j.expectedRes == "" {
					t.Errorf(`tcname : %s Expected result {""} but got {%s}`, j.tcName, val)
				} else {
					t.Errorf(`tcname : %s Expected result {%s} but got {%s}`, j.tcName, j.defaultVal, val)
				}
			}
		}
	}
}

func TestGetBool(t *testing.T) {
	testTable := []struct {
		tcName      string
		key         string
		value       string
		defaultVal  bool
		expectedRes any
	}{
		{"tc01", "UseCache", "false", true, false},                                  //passing tc
		{"tc02", "UseCache", "abc4d", false, "invalid value expected Boolean type"}, // should panic
		{"tc03", "UseCache", "", true, true},                                        //should take default value
		{"tc04", "UseCache", "true", false, true},                                   // shouldn't take default value
	}

	for _, j := range testTable {
		t.Run(j.tcName, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					err := fmt.Sprint(r)
					if err != j.expectedRes {
						t.Errorf("expected res: %s but got %s", j.expectedRes, err)
					}
				}
			}()

			if err := os.Setenv(j.key, j.value); err != nil {
				t.Error(err)
			}

			defer os.Unsetenv(j.key)
			val := GetBool(j.key, j.defaultVal)
			if j.key == "" {
				if j.defaultVal != val {
					t.Errorf("Expected result %t but got %t", j.defaultVal, val)
				}
			} else {
				if val != j.expectedRes {
					t.Errorf("Expected result %t but got %t", j.expectedRes, val)
				}
			}

		})
	}
}
func TestGetFloat64(t *testing.T) {
	testTable := []struct {
		tcName      string
		key         string
		value       string
		defaultVal  float64
		expectedRes any
	}{
		{"tc01", "MaxLoad", "3.54", 1, 3.54},                                       //passing tc
		{"tc02", "MaxLoad", "5", 7, 5.00},                                          // should panic
		{"tc03", "MaxLoad", "", 8.21, 8.21},                                        //should take default value
		{"tc04", "MaxLoad", "9.8", 3.5, 9.8},                                       // shouldn't take default value
		{"tc05", "MaxLoad", "no-load", 4.2, "invalid value expected float64 type"}, // should panic
	}

	for _, j := range testTable {
		t.Run(j.tcName, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					err := fmt.Sprint(r)
					if err != j.expectedRes {
						t.Errorf("expected res: %s but got %s", j.expectedRes, err)
					}
				}
			}()

			if err := os.Setenv(j.key, j.value); err != nil {
				t.Error(err)
			}

			defer os.Unsetenv(j.key)
			val := GetFloat64(j.key, j.defaultVal)
			if j.key == "" {
				if j.defaultVal != val {
					t.Errorf("Expected result %f but got %f", j.defaultVal, val)
				}
			} else {
				if val != j.expectedRes {
					t.Errorf("Expected result %f but got %f", j.expectedRes, val)
				}
			}

		})
	}
}

func TestGetStringArray(t *testing.T) {
	testTable := []struct {
		tcName      string
		key         string
		value       string
		defaultVal  []string
		expectedRes any
	}{
		{"tc01", "allowed_ip", "10.254.1.21,10.33.12.5,127.0.0.1", []string{"127.0.0.1,localhost:8080"}, []string{"10.254.1.21,10.33.12.5,127.0.0.1"}}, //passing tc
		{"tc02", "allowed_ip", "localhost:8080", []string{"12.22.7", "12.23.45.6"}, []string{"localhost:8080"}},                                        // should panic
		{"tc03", "allowed_ip", "", []string{"www.google.com", "www.mozila-firefox.com"}, []string{"www.google.com", "www.mozila-firefox.com"}},         //should take default value
	}

	for _, j := range testTable {
		t.Run(j.tcName, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					err := fmt.Sprint(r)
					if err != j.expectedRes {
						t.Errorf("expected res: %s but got %s", j.expectedRes, err)
					}
				}
			}()

			if err := os.Setenv(j.key, j.value); err != nil {
				t.Error(err)
			}

			defer os.Unsetenv(j.key)
			val := GetStringArray(j.key, j.defaultVal)
			if j.key == "" {
				if len(val) != len(j.defaultVal) {
					t.Errorf("invalid len of array expected %d but got %d", len(val), len(j.defaultVal))
				}
				for _, i := range val {
					found := false
					for _, j := range j.defaultVal {
						if i == j {
							found = true
							break
						}
					}
					if found == false {
						t.Errorf("invalid value of array expected %s couldn't be found", i)
					}
				}
			} else {
				for _, i := range val {
					found := false
					for _, j := range j.defaultVal {
						if i == j {
							found = true
							break
						}
					}
					if found == false {
						t.Errorf("invalid value of array expected %s couldn't be found", i)
					}
				}
			}

		})
	}
}

/*
func TestGetIntArray(t *testing.T) {
	testTable := []struct {
		tcName      string
		key         string
		value       string
		defaultVal  bool
		expectedRes any
	}{
		{"tc01", "Products", "false", true, false},                                  //passing tc
		{"tc02", "Products", "abc4d", false, "invalid value expected Boolean type"}, // should panic
		{"tc03", "Products", "", true, true},                                        //should take default value
		{"tc04", "Products", "true", false, true},                                   // shouldn't take default value
	}

	for _, j := range testTable {
		t.Run(j.tcName, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					err := fmt.Sprint(r)
					if err != j.expectedRes {
						t.Errorf("expected res: %s but got %s", j.expectedRes, err)
					}
				}
			}()

			if err := os.Setenv(j.key, j.value); err != nil {
				t.Error(err)
			}

			defer os.Unsetenv(j.key)
			val := GetBool(j.key, j.defaultVal)
			if j.key == "" {
				if j.defaultVal != val {
					t.Errorf("Expected result %t but got %t", j.defaultVal, val)
				}
			} else {
				if val != j.expectedRes {
					t.Errorf("Expected result %t but got %t", j.expectedRes, val)
				}
			}

		})
	}
}
*/
