package environment_test

import (
	"os"
	"testing"

	"github.com/blackskad/go-web-scaffold/pkg/environment"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	const envVarName = "GO-WEB-SCAFFOLD_ENABLE_TRACES_STDOUT"

	// Make sure the env var is not accidentially set somewhere else
	os.Unsetenv(envVarName)

	tests := []struct {
		name    string
		value   string
		enabled bool
	}{
		{
			name:    "default",
			value:   "",
			enabled: true,
		},
		{
			name:    "explicit true",
			value:   "true",
			enabled: true,
		},
		{
			name:    "explicit false",
			value:   "false",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(envVarName)

			if tt.value != "" {
				os.Setenv(envVarName, tt.value)
			}

			c := environment.Parse()
			assert.Equal(t, tt.enabled, c.EnableTracesStdout)
		})
	}
}
