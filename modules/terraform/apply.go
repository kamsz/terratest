package terraform

import (
	"errors"
	"testing"
)

// InitAndApply runs terraform init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func InitAndApply(t *testing.T, options *Options) string {
	out, err := InitAndApplyE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitAndApplyE runs terraform init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func InitAndApplyE(t *testing.T, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}

	if _, err := GetE(t, options); err != nil {
		return "", err
	}

	return ApplyE(t, options)
}

// Apply runs terraform apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func Apply(t *testing.T, options *Options) string {
	out, err := ApplyE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// ApplyAll runs terragrunt apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyAll(t *testing.T, options *Options) string {
	out, err := ApplyAllE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// ApplyE runs terraform apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyE(t *testing.T, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options.Vars, "apply", "-input=false", "-lock=false", "-auto-approve")...)
}

// ApplyAllE runs terragrunt apply-all with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyAllE(t *testing.T, options *Options) (string, error) {
	if options.TerraformBinary != "terragrunt" {
		return "", errors.New("terragrunt must be set as TerraformBinary to use this method")
	}

	return RunTerraformCommandE(t, options, FormatArgs(options.Vars, "apply-all", "-input=false", "-lock=false", "-auto-approve")...)
}
