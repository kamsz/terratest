package terraform

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/kamsz/terratest/modules/shell"
)

// RunTerraformCommand runs terraform with the given arguments and options and return stdout/stderr.
func RunTerraformCommand(t *testing.T, options *Options, args ...string) string {
	out, err := RunTerraformCommandE(t, options, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunTerraformCommandE runs terraform with the given arguments and options and return stdout/stderr.
func RunTerraformCommandE(t *testing.T, options *Options, args ...string) (string, error) {
	binary := options.TerraformBinary

	if binary == "" {
		binary = "terraform"
	}

	if binary == "terragrunt" {
		args = append(args, "--terragrunt-non-interactive")
	}

	if options.NoColor && !collections.ListContains(args, "-no-color") {
		args = append(args, "-no-color")
	}

	description := fmt.Sprintf("Running %s %v", binary, args)
	return retry.DoWithRetryE(t, description, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		cmd := shell.Command{
			Command:    binary,
			Args:       args,
			WorkingDir: options.TerraformDir,
			Env:        options.EnvVars,
			NoStderr:   options.NoStderr,
		}

		out, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil {
			return out, nil
		}

		for errorText, errorMessage := range options.RetryableTerraformErrors {
			if strings.Contains(out, errorText) {
				logger.Logf(t, "%s failed with the error '%s' but this error was expected and warrants a retry. Further details: %s\n", binary, errorText, errorMessage)
				return out, err
			}
		}

		return out, retry.FatalError{Underlying: err}
	})
}
