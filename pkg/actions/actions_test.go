package actions_test

import (
	"os"
	"testing"

	"github.com/gandarez/semver-action/pkg/actions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInput(t *testing.T) {
	tests := map[string]struct {
		Name        string
		NameActions string
		Expected    string
	}{
		"blank  spaces": {
			Name:        "Access Token",
			NameActions: "INPUT_ACCESS_TOKEN",
			Expected:    "my access token",
		},
		"all lower case": {
			Name:        "access_token",
			NameActions: "INPUT_ACCESS_TOKEN",
			Expected:    "my access token",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv(test.NameActions, test.Expected)
			defer os.Unsetenv(test.NameActions)

			value := actions.GetInput(test.Name)

			assert.Equal(t, test.Expected, value)
		})
	}
}

func TestSetOutput(t *testing.T) {
	outputFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	defer outputFile.Close()

	err = actions.SetOutput(outputFile.Name(), "SOME_OUTPUT", "some-value")
	require.NoError(t, err)

	data, err := os.ReadFile(outputFile.Name())
	require.NoError(t, err)

	assert.Regexp(
		t,
		`(?im)^.*<<ghadelimiter\_[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[0-9A-F]{4}-[0-9A-F]{12}\n.*\nghadelimiter\_[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[0-9A-F]{4}-[0-9A-F]{12}\n`,
		string(data),
	)
}
