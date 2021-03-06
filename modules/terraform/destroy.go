package terraform

import (
	"errors"
	"testing"
)

// Destroy runs terraform destroy with the given options and return stdout/stderr.
func Destroy(t *testing.T, options *Options) string {
	out, err := DestroyE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// DestroyAll runs terragrunt destroy with the given options and return stdout.
func DestroyAll(t *testing.T, options *Options) string {
	out, err := DestroyAllE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// DestroyE runs terraform destroy with the given options and return stdout/stderr.
func DestroyE(t *testing.T, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options.Vars, "destroy", "-force", "-input=false", "-lock=false")...)
}

// DestroyAllE runs terragrunt destroy with the given options and return stdout.
func DestroyAllE(t *testing.T, options *Options) (string, error) {
	if options.TerraformBinary != "terragrunt" {
		return "", errors.New("terragrunt must be set as TerraformBinary to use this method")
	}

	return RunTerraformCommandE(t, options, FormatArgs(options.Vars, "destroy-all", "-force", "-input=false", "-lock=false")...)
}
