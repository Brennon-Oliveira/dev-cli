package exec

import (
	"fmt"
	"strings"
)

type Executor interface {
	Run(name string, args ...string) error
	RunInteractive(name string, args ...string) error
	RunDetached(name string, args ...string) error
	Output(name string, args ...string) (string, error)
}

type RealExecutor struct{}

func NewExecutor() *RealExecutor {
	return &RealExecutor{}
}

type MockExecutor struct {
	Calls        []string
	OutputResult string
	OutputErr    error
	RunErr       error
}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		OutputResult: "",
		OutputErr:    nil,
		RunErr:       nil,
	}
}

func (m *MockExecutor) Run(name string, args ...string) error {
	m.Calls = append(m.Calls, formatCall(name, args...))
	return m.RunErr
}

func (m *MockExecutor) RunInteractive(name string, args ...string) error {
	m.Calls = append(m.Calls, formatCall(name, args...))
	return m.RunErr
}

func (m *MockExecutor) RunDetached(name string, args ...string) error {
	m.Calls = append(m.Calls, formatCall(name, args...))
	return m.RunErr
}

func (m *MockExecutor) Output(name string, args ...string) (string, error) {
	m.Calls = append(m.Calls, formatCall(name, args...))
	return m.OutputResult, m.OutputErr
}

func (m *MockExecutor) Reset() {
	m.Calls = nil
	m.OutputResult = ""
	m.OutputErr = nil
	m.RunErr = nil
}

func formatCall(name string, args ...string) string {
	if len(args) == 0 {
		return name
	}
	return fmt.Sprintf("%s %s", name, strings.Join(args, " "))
}
