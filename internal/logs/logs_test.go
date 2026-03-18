package logs

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "➜") {
		t.Errorf("expected output to contain info prefix, got: %s", output)
	}
}

func TestInfoWithArgs(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	Info("value: %s", "test")

	output := buf.String()
	if !strings.Contains(output, "value: test") {
		t.Errorf("expected output to contain 'value: test', got: %s", output)
	}
}

func TestVerbose_WhenDisabled(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetVerbose(false)

	Verbose("test message")

	if buf.Len() > 0 {
		t.Errorf("expected no output when verbose disabled, got: %s", buf.String())
	}
}

func TestVerbose_WhenEnabled(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetVerbose(true)

	Verbose("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "│") {
		t.Errorf("expected output to contain verbose prefix, got: %s", output)
	}
}

func TestVerboseWithArgs(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetVerbose(true)

	Verbose("executing: %s %s", "docker", "ps")

	output := buf.String()
	if !strings.Contains(output, "executing: docker ps") {
		t.Errorf("expected output to contain 'executing: docker ps', got: %s", output)
	}
}

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	Success("operation completed")

	output := buf.String()
	if !strings.Contains(output, "operation completed") {
		t.Errorf("expected output to contain 'operation completed', got: %s", output)
	}
	if !strings.Contains(output, "✓") {
		t.Errorf("expected output to contain success prefix, got: %s", output)
	}
}

func TestWarn(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	Warn("warning message")

	output := buf.String()
	if !strings.Contains(output, "warning message") {
		t.Errorf("expected output to contain 'warning message', got: %s", output)
	}
	if !strings.Contains(output, "⚠") {
		t.Errorf("expected output to contain warn prefix, got: %s", output)
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	Error("error occurred")

	output := buf.String()
	if !strings.Contains(output, "error occurred") {
		t.Errorf("expected output to contain 'error occurred', got: %s", output)
	}
	if !strings.Contains(output, "✗") {
		t.Errorf("expected output to contain error prefix, got: %s", output)
	}
}
