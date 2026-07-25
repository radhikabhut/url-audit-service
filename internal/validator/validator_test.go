package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSafetyURL(t *testing.T) {
	tests := []struct {
		name        string
		urlStr      string
		expectError bool
	}{
		{
			name:        "Valid HTTPS URL",
			urlStr:      "https://google.com",
			expectError: false,
		},
		{
			name:        "Valid HTTP URL",
			urlStr:      "http://example.com/some/path?param=1",
			expectError: false,
		},
		{
			name:        "Invalid scheme ftp",
			urlStr:      "ftp://google.com",
			expectError: true,
		},
		{
			name:        "Malformed URL",
			urlStr:      "://google.com",
			expectError: true,
		},
		{
			name:        "Localhost representation",
			urlStr:      "http://localhost",
			expectError: true,
		},
		{
			name:        "Localhost representation mixed-case",
			urlStr:      "http://LocalHost",
			expectError: true,
		},
		{
			name:        "Loopback IPv4 direct",
			urlStr:      "http://127.0.0.1",
			expectError: true,
		},
		{
			name:        "Loopback IPv6 direct",
			urlStr:      "http://[::1]",
			expectError: true,
		},
		{
			name:        "Private IPv4 Class A",
			urlStr:      "http://10.0.0.1",
			expectError: true,
		},
		{
			name:        "Private IPv4 Class C",
			urlStr:      "http://192.168.1.1",
			expectError: true,
		},
		{
			name:        "Unspecified IPv4",
			urlStr:      "http://0.0.0.0",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSafetyURL(tc.urlStr)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetValidatorAndValidateStruct(t *testing.T) {
	v := GetValidator()
	assert.NotNil(t, v)

	type TestStruct struct {
		Name string `validate:"required"`
	}

	// Valid
	err := ValidateStruct(TestStruct{Name: "Go"})
	assert.NoError(t, err)

	// Invalid
	err = ValidateStruct(TestStruct{})
	assert.Error(t, err)
}
