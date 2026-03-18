package container

import (
	"errors"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
)

func TestDockerClient_ListContainers(t *testing.T) {
	mockExec := exec.NewMockExecutor()
	client := NewDockerClient("docker", mockExec)

	err := client.ListContainers()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mockExec.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mockExec.Calls))
	}
}

func TestDockerClient_CleanResources(t *testing.T) {
	mockExec := exec.NewMockExecutor()
	client := NewDockerClient("docker", mockExec)

	err := client.CleanResources()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mockExec.Calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(mockExec.Calls))
	}
}

func TestDockerClient_CleanResources_Error(t *testing.T) {
	mockExec := exec.NewMockExecutor()
	mockExec.RunErr = errors.New("prune failed")
	client := NewDockerClient("docker", mockExec)

	err := client.CleanResources()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDockerClient_StopContainers(t *testing.T) {
	mockExec := exec.NewMockExecutor()
	client := NewDockerClient("docker", mockExec)

	err := client.StopContainers([]string{"id1", "id2"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mockExec.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mockExec.Calls))
	}
}

func TestDockerClient_StopContainers_Empty(t *testing.T) {
	mockExec := exec.NewMockExecutor()
	client := NewDockerClient("docker", mockExec)

	err := client.StopContainers([]string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mockExec.Calls) != 0 {
		t.Errorf("expected no calls for empty ids, got %d", len(mockExec.Calls))
	}
}

func TestDockerClient_RemoveContainers(t *testing.T) {
	mockExec := exec.NewMockExecutor()
	client := NewDockerClient("docker", mockExec)

	err := client.RemoveContainers([]string{"id1"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mockExec.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mockExec.Calls))
	}
}

func TestMockContainerClient(t *testing.T) {
	mock := NewMockContainerClient()

	_ = mock.ListContainers()
	_, _ = mock.GetContainerID("/test")
	_, _ = mock.GetAllRelatedContainers("/test")
	_ = mock.StopContainers([]string{"id1"})
	_ = mock.RemoveContainers([]string{"id1"})
	_ = mock.ShowLogs("/test", false)
	_ = mock.ListPorts("/test")
	_ = mock.CleanResources()

	if len(mock.Calls) != 8 {
		t.Errorf("expected 8 calls, got %d", len(mock.Calls))
	}
}
