package actions

import (
	"fmt"
	"os"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
)

// GetInput gets the input by the given name.
func GetInput(name string) string {
	e := strings.ReplaceAll(name, " ", "_")
	e = strings.ToUpper(e)
	e = "INPUT_" + e

	return strings.TrimSpace(os.Getenv(e))
}

// SetOutput sets the key value pair to output.
func SetOutput(fp, key, value string) error {
	f, err := os.OpenFile(fp, os.O_APPEND|os.O_WRONLY, 0600) // nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to open github output file: %s", err)
	}

	defer func() {
		_ = f.Close()
	}()

	id, err := newId()
	if err != nil {
		return err
	}

	delimiter := fmt.Sprintf("ghadelimiter_%s", id)

	if _, err := f.WriteString(fmt.Sprintf("%s<<%s\n%v\n%s\n", key, delimiter, value, delimiter)); err != nil {
		return fmt.Errorf("failed to write %s to output: %s", key, err)
	}

	return nil
}

func newId() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to created new uuid: %s", err)
	}

	return id.String(), nil
}
