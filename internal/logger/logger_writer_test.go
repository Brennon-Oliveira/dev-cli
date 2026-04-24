package logger

import (
	"bytes"
	"testing"

	loggerutils "github.com/Brennon-Oliveira/dev-cli/internal/logger/logger_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============ LoggerWriter Tests ============

func TestLoggerWriter_AllowWritingTrue(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	n, err := writer.Write([]byte("test data"))

	require.Nil(t, err)
	assert.Equal(t, 9, n) // "test data" is 9 bytes
	assert.NotEmpty(t, dest.String())
	assert.Contains(t, dest.String(), "test data")
}

func TestLoggerWriter_AllowWritingFalse(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, false)

	n, err := writer.Write([]byte("should not appear"))

	require.Nil(t, err)
	assert.Equal(t, 17, n) // Still returns length even if not written
	assert.Empty(t, dest.String())
}

func TestLoggerWriter_ReturnsCorrectLength(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, false)

	testData := []byte("hello world")
	n, err := writer.Write(testData)

	require.Nil(t, err)
	assert.Equal(t, len(testData), n)
}

func TestLoggerWriter_ReturnsCorrectLengthWhenWriting(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	testData := []byte("test data")
	n, err := writer.Write(testData)

	require.Nil(t, err)
	assert.Equal(t, len(testData), n)
}

func TestLoggerWriter_SetAllowWriting_Toggle(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, false)

	// Initially disabled
	writer.Write([]byte("first"))
	assert.Empty(t, dest.String())

	// Enable
	writer.SetAllowWriting(true)
	writer.Write([]byte("second"))
	assert.Contains(t, dest.String(), "second")

	// Disable
	dest.Reset()
	writer.SetAllowWriting(false)
	writer.Write([]byte("third"))
	assert.Empty(t, dest.String())
}

func TestLoggerWriter_FormattingWithPrefix(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	writer.Write([]byte("data"))

	output := dest.String()
	// Should contain the vertical bar prefix
	assert.Contains(t, output, "│")
	assert.Contains(t, output, loggerutils.HighIntensityBlackColor)
}

func TestLoggerWriter_FormattingWithSuffix(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	writer.Write([]byte("data"))

	output := dest.String()
	// Should contain reset color code at the end
	assert.Contains(t, output, loggerutils.ResetColor)
}

func TestLoggerWriter_MultipleWrites(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	writer.Write([]byte("first"))
	writer.Write([]byte("second"))
	writer.Write([]byte("third"))

	output := dest.String()
	assert.Contains(t, output, "first")
	assert.Contains(t, output, "second")
	assert.Contains(t, output, "third")
}

func TestLoggerWriter_EmptyData(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	n, err := writer.Write([]byte(""))

	require.Nil(t, err)
	assert.Equal(t, 0, n)
}

func TestLoggerWriter_LargeData(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = 'a'
	}

	n, err := writer.Write(largeData)

	require.Nil(t, err)
	assert.Equal(t, 10000, n)
	assert.Contains(t, dest.String(), string(largeData))
}

func TestLoggerWriter_BinaryData(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	binaryData := []byte{0x00, 0x01, 0x02, 0xFF}
	n, err := writer.Write(binaryData)

	require.Nil(t, err)
	assert.Equal(t, 4, n)
}

func TestLoggerWriter_StateIsIndependent(t *testing.T) {
	dest1 := new(bytes.Buffer)
	dest2 := new(bytes.Buffer)

	writer1 := NewLoggerWriter(dest1, true)
	writer2 := NewLoggerWriter(dest2, false)

	writer1.Write([]byte("data1"))
	writer2.Write([]byte("data2"))

	assert.Contains(t, dest1.String(), "data1")
	assert.Empty(t, dest2.String())
}

func TestLoggerWriter_ToggleIndependentPerInstance(t *testing.T) {
	dest1 := new(bytes.Buffer)
	dest2 := new(bytes.Buffer)

	writer1 := NewLoggerWriter(dest1, true)
	writer2 := NewLoggerWriter(dest2, true)

	// Disable one
	writer1.SetAllowWriting(false)

	writer1.Write([]byte("disabled"))
	writer2.Write([]byte("enabled"))

	assert.Empty(t, dest1.String())
	assert.Contains(t, dest2.String(), "enabled")
}

func TestLoggerWriter_SpecialCharacters(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	specialChars := []byte("special: @#$%^&*()")
	n, err := writer.Write(specialChars)

	require.Nil(t, err)
	assert.Equal(t, len(specialChars), n)
	assert.Contains(t, dest.String(), "special: @#$%^&*()")
}

func TestLoggerWriter_UnicodeCharacters(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	unicodeData := []byte("unicode: ✓ ✗ ⚠ ➜")
	n, err := writer.Write(unicodeData)

	require.Nil(t, err)
	assert.Equal(t, len(unicodeData), n)
	assert.Contains(t, dest.String(), "unicode: ✓ ✗ ⚠ ➜")
}

func TestLoggerWriter_Newlines(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	dataWithNewlines := []byte("line1\nline2\nline3")
	n, err := writer.Write(dataWithNewlines)

	require.Nil(t, err)
	assert.Equal(t, len(dataWithNewlines), n)
	assert.Contains(t, dest.String(), "line1")
	assert.Contains(t, dest.String(), "line2")
	assert.Contains(t, dest.String(), "line3")
}

func TestLoggerWriter_ExclamationMarks(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	writer.Write([]byte("!!!"))

	assert.Contains(t, dest.String(), "!!!")
}

func TestLoggerWriter_Tabs(t *testing.T) {
	dest := new(bytes.Buffer)
	writer := NewLoggerWriter(dest, true)

	writer.Write([]byte("col1\tcol2\tcol3"))

	assert.Contains(t, dest.String(), "col1\tcol2\tcol3")
}
