package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// Command is a simpler struct for defining commands than Go's built-in Cmd.
type Command struct {
	Command    string            // The command to run
	Args       []string          // The args to pass to the command
	WorkingDir string            // The working directory
	Env        map[string]string // Additional environment variables to set
	NoStderr   bool              // Redirect stderr to output
}

// RunCommand runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself.
func RunCommand(t *testing.T, command Command) {
	err := RunCommandE(t, command)
	if err != nil {
		t.Fatal(err)
	}
}

// RunCommandE runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself.
func RunCommandE(t *testing.T, command Command) error {
	_, err := RunCommandAndGetOutputE(t, command)
	return err
}

// RunCommandAndGetOutput runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of that command will also
// be printed to the stdout and stderr of this Go program to make debugging easier.
func RunCommandAndGetOutput(t *testing.T, command Command) string {
	out, err := RunCommandAndGetOutputE(t, command)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunCommandAndGetOutputE runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of that command will also
// be printed to the stdout and stderr of this Go program to make debugging easier.
func RunCommandAndGetOutputE(t *testing.T, command Command) (string, error) {
	logger.Logf(t, "Running command %s with args %s", command.Command, command.Args)

	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Env = formatEnvVars(command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	var output string
	if command.NoStderr {
		output, err = readStdoutAndStderr(t, stdout)
	} else {
		output, err = readStdoutAndStderr(t, stdout, stderr)
	}

	if err != nil {
		return output, err
	}

	if err := cmd.Wait(); err != nil {
		return output, err
	}

	return output, nil
}

// This function captures stdout and stderr while still printing it to the stdout and stderr of this Go program
func readStdoutAndStderr(t *testing.T, outputs ...io.ReadCloser) (string, error) {
	allOutput := []string{}

	scanners := []bufio.Scanner{}
	for i := range outputs {
		scanners = append(scanners, *bufio.NewScanner(outputs[i]))
	}

	for i := range scanners {
		if scanners[i].Scan() {
			text := scanners[i].Text()
			logger.Log(t, text)
			allOutput = append(allOutput, text)
		}

		if err := scanners[i].Err(); err != nil {
			return "", err
		}
	}

	return strings.Join(allOutput, "\n"), nil
}

// GetExitCodeForRunCommandError tries to read the exit code for the error object returned from running a shell command. This is a bit tricky to do
// in a way that works across platforms.
func GetExitCodeForRunCommandError(err error) (int, error) {
	// http://stackoverflow.com/a/10385867/483528
	if exitErr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), nil
		}
		return 1, errors.New("could not determine exit code")
	}

	return 0, nil
}

func formatEnvVars(command Command) []string {
	env := os.Environ()
	for key, value := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}
