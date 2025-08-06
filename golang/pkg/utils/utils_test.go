package utils

import "testing"

func TestGetEnv(t *testing.T) {
	value := GetEnv("TEST_ENV", "default_value")
	if value != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", value)
	}

	// Set an environment variable for testing
	t.Setenv("TEST_ENV", "test_value")
	value = GetEnv("TEST_ENV", "default_value")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}
func TestRequireNotEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for empty value, but did not panic")
		}
	}()

	RequireNotEmpty("", "testValue")

	t.Setenv("TEST_ENV", "non_empty_value")
	RequireNotEmpty(GetEnv("TEST_ENV", ""), "TEST_ENV")
}
