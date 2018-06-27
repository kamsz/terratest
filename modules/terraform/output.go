package terraform

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Output calls terraform output for the given variable and return its value.
func Output(t *testing.T, options *Options, key string) string {
	out, err := OutputE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// OutputAll calls terragrunt output for the given variable and return its value.
func OutputAll(t *testing.T, options *Options, key string) string {
	out, err := OutputAllE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// OutputE calls terraform output for the given variable and return its value.
func OutputE(t *testing.T, options *Options, key string) (string, error) {
	if options.TerraformBinary == "terragrunt" {
		options.NoStderr = true
	}

	output, err := RunTerraformCommandE(t, options, "output", "-no-color", key)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// OutputAllE calls terragrunt output for the given variable and return its value.
func OutputAllE(t *testing.T, options *Options, key string) (string, error) {
	if options.TerraformBinary != "terragrunt" {
		return "", errors.New("terragrunt must be set as TerraformBinary to use this method")
	}

	options.NoStderr = true
	output, err := RunTerraformCommandE(t, options, "output-all", "-no-color", key)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// OutputRequired calls terraform output for the given variable and return its value. If the value is empty, fail the test.
func OutputRequired(t *testing.T, options *Options, key string) string {
	out, err := OutputRequiredE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// OutputRequiredE calls terraform output for the given variable and return its value. If the value is empty, return an error.
func OutputRequiredE(t *testing.T, options *Options, key string) (string, error) {
	out, err := OutputE(t, options, key)

	if err != nil {
		return "", err
	}
	if out == "" {
		return "", EmptyOutput(key)
	}

	return out, nil
}

// EmptyOutput is an error that occurs when an output is empty.
type EmptyOutput string

func (outputName EmptyOutput) Error() string {
	return fmt.Sprintf("Required output %s was empty", string(outputName))
}
