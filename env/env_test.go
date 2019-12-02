package env

import (
	"os"
	"testing"
)

func TestSetDefaultEnvValue(t *testing.T) {

	const Default = "testDefaultValue"

	// when an environment variable is not set use default
	os.Unsetenv("TEST_ENV")
	value := GetWithDefault("TEST_ENV", Default)
	if value != Default {
		t.Errorf(`expected value to be "%s" (default), got "%s"`, Default, value)
	}
}

func TestSetDefaultEnvValue2(t *testing.T) {

	const Default = "testDefaultValue"

	// when an environment variable is set use return its value
	os.Setenv("TEST_ENV", "test_value")
	value := GetWithDefault("TEST_ENV", Default)
	if value != "test_value" {
		t.Errorf(`expected value to be "%s" (environment variable value), got "%s"`, "test_value", value)
	}
}
