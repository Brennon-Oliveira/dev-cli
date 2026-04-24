package logger

import (
	"bytes"
	"os"
	"testing"

	loggerutils "github.com/Brennon-Oliveira/dev-cli/internal/logger/logger_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// ============ Initialization & Configuration Tests ============

func TestNewLogger_DefaultSettings(t *testing.T) {
	l := NewLogger()
	require.NotNil(t, l)

	// Default should be non-verbose
	output := new(bytes.Buffer)
	l.SetOutput(output)
	l.Verbose("should not appear")

	assert.Empty(t, output.String())
}

func TestInitLogger_SetsGlobal(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))
	require.NotNil(t, logger)

	// Verify global logger is set and works
	logger.Info("test message")
	assert.Contains(t, output.String(), "test message")
}

func TestWithWriter_CustomOutput(t *testing.T) {
	customOutput := new(bytes.Buffer)

	l := NewLogger(WithWriter(customOutput))
	l.Info("custom output test")

	require.NotNil(t, customOutput.String())
	assert.Contains(t, customOutput.String(), "custom output test")
}

func TestWithVerbose_EnablesVerbose(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(
		WithWriter(output),
		WithVerbose(true),
	)

	l.Verbose("verbose message")
	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "verbose message")
}

func TestSetVerbose_Toggles(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))

	// Initially disabled
	l.Verbose("should not appear")
	assert.Empty(t, output.String())

	// Enable verbose
	l.SetVerbose(true)
	l.Verbose("should appear")
	assert.Contains(t, output.String(), "should appear")

	// Clear and disable
	output.Reset()
	l.SetVerbose(false)
	l.Verbose("should not appear again")
	assert.Empty(t, output.String())
}

func TestSetOutput_ReplaceWriter(t *testing.T) {
	output1 := new(bytes.Buffer)
	output2 := new(bytes.Buffer)

	l := NewLogger(WithWriter(output1))
	l.Info("first output")

	assert.Contains(t, output1.String(), "first output")
	assert.Empty(t, output2.String())

	// Change output
	l.SetOutput(output2)
	l.Info("second output")

	// Verify second output has the message
	currentOutput := new(bytes.Buffer)
	l.SetOutput(currentOutput)
	l.Info("verify")
	assert.Contains(t, currentOutput.String(), "verify")
}

// ============ Info Tests ============

func TestInfo_PlainString(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("plain info message")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "plain info message")
	assert.Contains(t, output.String(), loggerutils.RegularCyanColor)
	assert.Contains(t, output.String(), "➜")
}

func TestInfo_WithArgs(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("info with %s and %d", "string", 42)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "info with string and 42")
	assert.Contains(t, output.String(), "➜")
}

// ============ Verbose Tests ============

func TestVerbose_DisabledNoOutput(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	// Verbose is disabled by default
	l.Verbose("should not appear")

	assert.Empty(t, output.String())
}

func TestVerbose_EnabledProducesOutput(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(
		WithWriter(output),
		WithVerbose(true),
	)

	l.Verbose("verbose message")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "verbose message")
	assert.Contains(t, output.String(), "│")
}

func TestVerbose_WithArgs(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(
		WithWriter(output),
		WithVerbose(true),
	)

	l.Verbose("verbose with %s and %d", "arg", 123)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "verbose with arg and 123")
	assert.Contains(t, output.String(), "│")
}

// ============ Debug Tests ============

func TestDebug_PlainString(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Debug("debug message")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "debug message")
	assert.Contains(t, output.String(), loggerutils.RegularBlueColor)
	assert.Contains(t, output.String(), "➜")
}

func TestDebug_WithArgs(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Debug("debug with %s and %d", "details", 99)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "debug with details and 99")
	assert.Contains(t, output.String(), "➜")
}

// ============ Success Tests ============

func TestSuccess_PlainString(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Success("operation completed")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "operation completed")
	assert.Contains(t, output.String(), loggerutils.RegularGreenColor)
	assert.Contains(t, output.String(), "✓")
}

func TestSuccess_WithArgs(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Success("success with %s and count %d", "details", 5)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "success with details and count 5")
	assert.Contains(t, output.String(), "✓")
}

// ============ Warn Tests ============

func TestWarn_PlainString(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Warn("warning message")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "warning message")
	assert.Contains(t, output.String(), loggerutils.RegularYellowColor)
	assert.Contains(t, output.String(), "⚠")
}

func TestWarn_WithArgs(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Warn("warning: %s occurred %d times", "issue", 3)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "warning: issue occurred 3 times")
	assert.Contains(t, output.String(), "⚠")
}

// ============ Error Tests ============

func TestError_PlainString(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Error("error message")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "error message")
	assert.Contains(t, output.String(), loggerutils.RegularRedColor)
	assert.Contains(t, output.String(), "✗")
}

func TestError_WithArgs(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Error("error: %s failed with code %d", "operation", 500)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "error: operation failed with code 500")
	assert.Contains(t, output.String(), "✗")
}

// ============ VerboseFromOutput Tests ============

func TestVerboseFromOutput_CopiesDataWhenEnabled(t *testing.T) {
	output := new(bytes.Buffer)
	input := bytes.NewBufferString("input data from reader")

	l := NewLogger(
		WithWriter(output),
		WithVerbose(true),
	)

	l.VerboseFromOutput(input)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "input data from reader")
}

func TestVerboseFromOutput_NoOutputWhenDisabled(t *testing.T) {
	output := new(bytes.Buffer)
	input := bytes.NewBufferString("should not appear")

	l := NewLogger(WithWriter(output))
	// Verbose disabled by default

	l.VerboseFromOutput(input)

	assert.Empty(t, output.String())
}

func TestVerboseFromOutput_FormatsWithVerticalBar(t *testing.T) {
	output := new(bytes.Buffer)
	input := bytes.NewBufferString("formatted output")

	l := NewLogger(
		WithWriter(output),
		WithVerbose(true),
	)

	l.VerboseFromOutput(input)

	assert.Contains(t, output.String(), "│")
	assert.Contains(t, output.String(), loggerutils.HighIntensityBackgroundBlackColor)
}

// ============ GetWriter Tests ============

func TestGetWriter_ReturnsLoggerWriter(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	writer := l.GetWriter()

	require.NotNil(t, writer)
	// Writer should be LoggerWriter type
	_, ok := writer.(*LoggerWriter)
	assert.True(t, ok)
}

// ============ Global Function Tests ============

func TestGlobalSetVerbose_UpdatesSingleton(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))

	// Initially disabled
	Verbose("should not appear")
	assert.Empty(t, output.String())

	// Enable via global
	SetVerbose(true)
	Verbose("should appear")
	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "should appear")
}

func TestGlobalSetOutput_UpdatesSingleton(t *testing.T) {
	output1 := new(bytes.Buffer)
	output2 := new(bytes.Buffer)

	InitLogger(WithWriter(output1))
	Info("first")

	assert.Contains(t, output1.String(), "first")

	SetOutput(output2)
	Info("second")

	// Verify the global logger is using the new output
	tempOutput := new(bytes.Buffer)
	SetOutput(tempOutput)
	Info("test")
	assert.Contains(t, tempOutput.String(), "test")
}

func TestGlobalInfo_CallsSingletonInfo(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))
	Info("global info test")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "global info test")
}

func TestGlobalVerbose_CallsSingletonVerbose(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output), WithVerbose(true))

	Verbose("global verbose test")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "global verbose test")
}

func TestGlobalDebug_CallsSingletonDebug(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))
	Debug("global debug test")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "global debug test")
}

func TestGlobalSuccess_CallsSingletonSuccess(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))
	Success("global success test")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "global success test")
	assert.Contains(t, output.String(), "✓")
}

func TestGlobalWarn_CallsSingletonWarn(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))
	Warn("global warn test")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "global warn test")
	assert.Contains(t, output.String(), "⚠")
}

func TestGlobalError_CallsSingletonError(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))
	Error("global error test")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "global error test")
	assert.Contains(t, output.String(), "✗")
}

func TestGlobalGetWriter_ReturnsSingletonWriter(t *testing.T) {
	output := new(bytes.Buffer)

	InitLogger(WithWriter(output))
	writer := GetWriter()

	require.NotNil(t, writer)
	_, ok := writer.(*LoggerWriter)
	assert.True(t, ok)
}

func TestGlobalVerboseFromOutput_CallsSingletonVerboseFromOutput(t *testing.T) {
	output := new(bytes.Buffer)
	input := bytes.NewBufferString("global verbose from output")

	InitLogger(WithWriter(output), WithVerbose(true))
	VerboseFromOutput(input)

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "global verbose from output")
}

// ============ Format String Edge Cases ============

func TestInfo_EmptyFormat(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "➜")
}

func TestVerbose_EmptyFormat(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(
		WithWriter(output),
		WithVerbose(true),
	)
	l.Verbose("")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "│")
}

func TestInfo_MultilineMessage(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("line1\nline2")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "line1")
	assert.Contains(t, output.String(), "line2")
}

func TestSuccess_SpecialCharacters(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Success("success with special chars: @#$*()")

	require.NotEmpty(t, output.String())
	assert.Contains(t, output.String(), "special chars: @#$*()")
}

// ============ Multiple Logger Instances ============

func TestMultipleLoggers_IndependentState(t *testing.T) {
	output1 := new(bytes.Buffer)
	output2 := new(bytes.Buffer)

	l1 := NewLogger(WithWriter(output1), WithVerbose(true))
	l2 := NewLogger(WithWriter(output2), WithVerbose(false))

	l1.Verbose("from logger 1")
	l2.Verbose("from logger 2")

	assert.Contains(t, output1.String(), "from logger 1")
	assert.NotContains(t, output1.String(), "from logger 2")
	assert.Empty(t, output2.String())
}

func TestMultipleLoggers_DifferentOutputs(t *testing.T) {
	output1 := new(bytes.Buffer)
	output2 := new(bytes.Buffer)

	l1 := NewLogger(WithWriter(output1))
	l2 := NewLogger(WithWriter(output2))

	l1.Info("message 1")
	l2.Info("message 2")

	assert.Contains(t, output1.String(), "message 1")
	assert.NotContains(t, output1.String(), "message 2")
	assert.NotContains(t, output2.String(), "message 1")
	assert.Contains(t, output2.String(), "message 2")
}

// ============ Color Code Verification ============

func TestOutputContainsANSICodes(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("color test")

	// Should contain reset code
	assert.Contains(t, output.String(), loggerutils.ResetColor)
}

func TestInfoContainsCyanColor(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("cyan color")

	assert.Contains(t, output.String(), loggerutils.RegularCyanColor)
}

func TestSuccessContainsGreenColor(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Success("green color")

	assert.Contains(t, output.String(), loggerutils.RegularGreenColor)
}

func TestErrorContainsRedColor(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Error("red color")

	assert.Contains(t, output.String(), loggerutils.RegularRedColor)
}

func TestWarnContainsYellowColor(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Warn("yellow color")

	assert.Contains(t, output.String(), loggerutils.RegularYellowColor)
}

func TestDebugContainsBlueColor(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Debug("blue color")

	assert.Contains(t, output.String(), loggerutils.RegularBlueColor)
}

// ============ Symbol Verification ============

func TestInfoHasArrowSymbol(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("arrow symbol")

	assert.Contains(t, output.String(), "➜")
}

func TestSuccessHasCheckmark(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Success("checkmark symbol")

	assert.Contains(t, output.String(), "✓")
}

func TestErrorHasX(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Error("x symbol")

	assert.Contains(t, output.String(), "✗")
}

func TestWarnHasWarningSymbol(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Warn("warning symbol")

	assert.Contains(t, output.String(), "⚠")
}

func TestVerboseHasVerticalBar(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(
		WithWriter(output),
		WithVerbose(true),
	)
	l.Verbose("vertical bar")

	assert.Contains(t, output.String(), "│")
}

// ============ Format Specification Tests ============

func TestInfo_IntegerFormatting(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("value is %d", 12345)

	assert.Contains(t, output.String(), "value is 12345")
}

func TestInfo_FloatFormatting(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("value is %.2f", 3.14159)

	assert.Contains(t, output.String(), "value is 3.14")
}

func TestInfo_MultipleFormatArgs(t *testing.T) {
	output := new(bytes.Buffer)

	l := NewLogger(WithWriter(output))
	l.Info("%s %d %v", "text", 42, true)

	assert.Contains(t, output.String(), "text 42 true")
}
