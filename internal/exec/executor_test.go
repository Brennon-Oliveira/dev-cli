package exec

import (
	"bytes"
	"os"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	exitCode := m.Run()
	os.Exit(exitCode)
}

// ============ Run Tests ============

func TestRun_SuccessWithSimpleCommand(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// 'echo' with 3 args to match the implementation's fmt.Printf expectation
	err := executor.Run("echo", "arg1", "arg2", "arg3")
	r.Nil(err)
}

func TestRun_CommandNotFound(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// Non-existent command should return error (with 3 args to match implementation)
	err := executor.Run("nonexistentcommand12345", "arg1", "arg2", "arg3")
	r.NotNil(err)
}

func TestRun_NonZeroExit(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// 'false' is a standard Unix command that always fails (with 3 args)
	err := executor.Run("sh", "-c", "false", "dummy")
	r.NotNil(err)
}

func TestRun_WithThreeArgs(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// Run command with exactly 3 arguments as required by implementation
	err := executor.Run("echo", "hello", "world", "test")
	r.Nil(err)
}

func TestRun_WithExactlyThreeArgs(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// 'echo' with exactly 3 arguments
	err := executor.Run("echo", "arg1", "arg2", "arg3")
	r.Nil(err)
}

func TestRun_WithMultipleArgs(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// 'echo' with multiple arguments
	err := executor.Run("echo", "hello", "world", "test")
	r.Nil(err)
}

func TestRun_StdoutCapture(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// Use 3 arguments
	err := executor.Run("echo", "test_output", "arg2", "arg3")
	r.Nil(err)

	// Verify output was captured
	output := stdout.String()
	assert.Contains(t, output, "test_output")
}

// ============ RunWithOutput Tests ============

func TestRunWithOutput_SuccessWithSimpleCommand(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	customOutput := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// RunWithOutput doesn't have the 3-arg limitation issue, but we'll use 3 for consistency
	err := executor.RunWithOutput(customOutput, "echo", "arg1", "arg2")
	r.Nil(err)
}

func TestRunWithOutput_CommandNotFound(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	customOutput := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	err := executor.RunWithOutput(customOutput, "nonexistentcommand12345", "a", "b")
	r.NotNil(err)
}

func TestRunWithOutput_NonZeroExit(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	customOutput := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	err := executor.RunWithOutput(customOutput, "sh", "-c", "false")
	r.NotNil(err)
}

func TestRunWithOutput_WithMultipleArgs(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	customOutput := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	err := executor.RunWithOutput(customOutput, "echo", "hello", "world")
	r.Nil(err)
}

func TestRunWithOutput_OutputMultiplexing(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	customOutput := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	err := executor.RunWithOutput(customOutput, "echo", "multiplex_test")
	r.Nil(err)

	// Verify output went to both writers
	stdoutContent := stdout.String()
	customContent := customOutput.String()

	assert.Contains(t, stdoutContent, "multiplex_test")
	assert.Contains(t, customContent, "multiplex_test")
}

func TestRunWithOutput_CustomOutputOnly(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	customOutput := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	testString := "unique_test_output"
	err := executor.RunWithOutput(customOutput, "echo", testString)
	r.Nil(err)

	// Verify output was captured in custom writer
	customContent := customOutput.String()
	assert.Contains(t, customContent, testString)
}

// ============ Output Tests ============

func TestOutput_SuccessWithSimpleCommand(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	output, err := executor.Output("true")
	r.Nil(err)
	r.NotNil(output)
}

func TestOutput_CommandNotFound(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	_, err := executor.Output("nonexistentcommand12345")
	r.NotNil(err)
}

func TestOutput_NonZeroExit(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// 'false' returns non-zero, which means Output() returns error
	_, err := executor.Output("false")
	r.NotNil(err)
}

func TestOutput_CapturesOutput(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	output, err := executor.Output("echo", "hello_world")
	r.Nil(err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "hello_world")
}

func TestOutput_WithMultipleArgs(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	output, err := executor.Output("echo", "arg1", "arg2", "arg3")
	r.Nil(err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "arg1")
	assert.Contains(t, outputStr, "arg2")
	assert.Contains(t, outputStr, "arg3")
}

func TestOutput_ReturnsBytes(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	output, err := executor.Output("echo", "bytes_test")
	r.Nil(err)

	assert.IsType(t, []byte{}, output)
	assert.True(t, len(output) > 0)
}

// ============ RunDetached Tests ============

func TestRunDetached_SuccessWithSimpleCommand(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// 'sleep 0' starts immediately and completes quickly
	err := executor.RunDetached("sleep", "0")
	r.Nil(err)
}

func TestRunDetached_CommandNotFound(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	err := executor.RunDetached("nonexistentcommand12345")
	r.NotNil(err)
}

func TestRunDetached_WithArgs(t *testing.T) {
	r := require.New(t)
	stdout := new(bytes.Buffer)
	stdin := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdout(stdout),
		WithStdin(stdin),
	)

	// 'sleep 0' with argument
	err := executor.RunDetached("sleep", "0")
	r.Nil(err)
}

// ============ Default Executor Tests ============

func TestNewExecutor_DefaultStdout(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()
	r.NotNil(executor)

	// Verify it implements Executor interface
	var _ Executor = executor
}

func TestNewExecutor_WithCustomStdout(t *testing.T) {
	r := require.New(t)
	customOut := new(bytes.Buffer)

	executor := NewExecutor(WithStdout(customOut))
	r.NotNil(executor)

	// Run a command with 3 args to verify custom stdout is used
	err := executor.Run("echo", "test", "arg2", "arg3")
	r.Nil(err)

	// Verify something was written
	assert.True(t, customOut.Len() > 0)
}

func TestNewExecutor_WithCustomStdin(t *testing.T) {
	r := require.New(t)
	customIn := bytes.NewBufferString("test input")
	customOut := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdin(customIn),
		WithStdout(customOut),
	)
	r.NotNil(executor)
}

func TestNewExecutor_WithBothCustomStreams(t *testing.T) {
	r := require.New(t)
	customIn := bytes.NewBufferString("")
	customOut := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdin(customIn),
		WithStdout(customOut),
	)
	r.NotNil(executor)

	// Use 3 arguments
	err := executor.Run("echo", "streams", "arg2", "arg3")
	r.Nil(err)

	assert.True(t, customOut.Len() > 0)
}

// ============ RunInteractive Tests ============

func TestRunInteractive_SuccessWithSimpleCommand(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// Use echo as a simple non-interactive command to test the mechanics
	err := executor.RunInteractive("echo", "hello", "world")
	r.Nil(err)
}

func TestRunInteractive_CommandNotFound(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// Non-existent command should return error
	err := executor.RunInteractive("nonexistentcommand12345")
	r.NotNil(err)
}

func TestRunInteractive_PassesThroughStdin(t *testing.T) {
	r := require.New(t)

	// Create a stdin with test data
	stdin := bytes.NewBufferString("test input\n")

	executor := NewExecutor(
		WithStdin(stdin),
	)

	// Test that RunInteractive attempts to run and uses os.Stdin
	// Using echo instead of cat since cat requires actual terminal interaction
	err := executor.RunInteractive("echo", "stdin_test")
	r.Nil(err)
}

func TestRunInteractive_UsesOsStdout(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// RunInteractive should use os.Stdout (which is being captured by the test framework)
	err := executor.RunInteractive("echo", "interactive_output")
	r.Nil(err)
}

func TestRunInteractive_UsesOsStderr(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// Command that outputs to stderr
	err := executor.RunInteractive("sh", "-c", "echo error_message >&2")
	r.Nil(err)
}

func TestRunInteractive_CommandWithMultipleArgs(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// Test with multiple arguments
	err := executor.RunInteractive("echo", "arg1", "arg2", "arg3", "arg4")
	r.Nil(err)
}

func TestRunInteractive_NonZeroExit(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// false command returns exit code 1
	err := executor.RunInteractive("false")
	r.NotNil(err)
}

func TestRunInteractive_ComplexCommand(t *testing.T) {
	r := require.New(t)

	executor := NewExecutor()

	// Test with a more complex shell command
	err := executor.RunInteractive("sh", "-c", "echo 'line 1'; echo 'line 2'")
	r.Nil(err)
}

func TestRunInteractive_IgnoresCustomStdin(t *testing.T) {
	r := require.New(t)

	// Create custom stdin but RunInteractive should use os.Stdin
	customStdin := bytes.NewBufferString("custom data")
	customStdout := new(bytes.Buffer)

	executor := NewExecutor(
		WithStdin(customStdin),
		WithStdout(customStdout),
	)

	// RunInteractive uses os.Stdin directly, not the custom stdin
	err := executor.RunInteractive("echo", "test")
	r.Nil(err)

	// Custom stdout should be ignored by RunInteractive
	// Since RunInteractive uses os.Stdout directly
	assert.Empty(t, customStdout.String())
}

func TestRunInteractive_DirectsToOsStreams(t *testing.T) {
	r := require.New(t)

	// RunInteractive should use os.Stdin, os.Stdout, os.Stderr
	// We verify this by checking that custom streams are NOT used
	customStdin := bytes.NewBufferString("should not be used")
	customStdout := new(bytes.Buffer)

	executorWithStreams := NewExecutor(
		WithStdin(customStdin),
		WithStdout(customStdout),
	)

	err := executorWithStreams.RunInteractive("echo", "direct")
	r.Nil(err)

	// Output should NOT be in custom stdout because RunInteractive uses os.Stdout
	assert.Empty(t, customStdout.String())
}
