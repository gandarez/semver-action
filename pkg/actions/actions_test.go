package actions_test

import (
	"os"
	"testing"

	"github.com/gandarez/semver-action/pkg/actions"
	"github.com/stretchr/testify/assert"
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
