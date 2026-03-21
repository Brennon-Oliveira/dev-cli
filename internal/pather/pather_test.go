package pather

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetAbsPath_EmptyReturnsCWD(t *testing.T) {
	r := require.New(t)
	cwd, err := os.Getwd()
	r.Nil(err)

	pather := NewPather()

	got, err := pather.GetAbsPath("")
	r.Nil(err)
	assert.Equal(t, cwd, got)
}

func TestGetAbsPath_DotReturnsCWD(t *testing.T) {
	r := require.New(t)
	cwd, err := os.Getwd()
	r.Nil(err)

	pather := NewPather()
	got, err := pather.GetAbsPath(".")
	r.Nil(err)

	assert.Equal(t, cwd, got)
}

func TestGetAbsPath_RelativePath(t *testing.T) {
	r := require.New(t)
	cwd, err := os.Getwd()
	r.Nil(err)

	pather := NewPather()

	got, err := pather.GetAbsPath("subdir")
	r.Nil(err)

	expected := filepath.Join(cwd, "subdir")
	assert.Equal(t, expected, got)
}

func TestGetAbsPath_AbsolutePathUncharged(t *testing.T) {
	r := require.New(t)
	absPath := "/tmp/test"
	pather := NewPather()
	got, err := pather.GetAbsPath(absPath)
	r.Nil(err)
	assert.Equal(t, absPath, got)
}

func TestGetPathFromArgs_Empty(t *testing.T) {
	args := []string{}
	pather := NewPather()
	got := pather.GetPathFromArgs(args)
	assert.Empty(t, got)
}

func TestGetPathFromArgs_RelativeWithDotPath(t *testing.T) {
	path := "./subdir"
	args := []string{path}
	pather := NewPather()
	got := pather.GetPathFromArgs(args)
	assert.Equal(t, path, got)
}

func TestGetPathFromArgs_RelativeWithoutDotPath(t *testing.T) {
	path := "subdir"
	args := []string{path}
	pather := NewPather()
	got := pather.GetPathFromArgs(args)
	assert.Equal(t, path, got)
}

func TestGetPathFromArgs_AbsolutePath(t *testing.T) {
	path := "/tmp/test"
	args := []string{path}
	pather := NewPather()
	got := pather.GetPathFromArgs(args)
	assert.Equal(t, path, got)
}

func lookupEnvInWslMock(key string) (string, bool) {
	if key == "WSL_DISTRO_NAME" {
		return "Ubuntu", true
	}
	return "", false
}

func TestGetRealPath_WSLReturnsCorrect(t *testing.T) {
	r := require.New(t)
	wslPath := "/tmp/wsl"
	t.Setenv("WSL_DISTRO_NAME", "Ubuntu")
	path := "/tmp/my-project"
	fullPath := wslPath + path

	executor := exec.NewMockExecutor(t)

	executor.EXPECT().RunWithOutput(mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(output io.Writer, name string, args ...string) {
		fmt.Fprint(output, fullPath)
	})

	pather := NewPather(
		WithExecutor(executor),
		WithLookupEnv(lookupEnvInWslMock),
	)

	got, err := pather.GetRealPath(path)
	r.Nil(err)
	assert.Equal(t, fullPath, got)
}

func lookupEnvOutWslMock(key string) (string, bool) {
	if key == "WSL_DISTRO_NAME" {
		return "", false
	}
	return "", true
}

func TestGetRealPath_NoWSLReturnsOriginal(t *testing.T) {
	r := require.New(t)
	t.Setenv("WSL_DISTRO_NAME", "")
	path := "/tmp/my-project"

	executor := exec.NewMockExecutor(t)

	executor.AssertNotCalled(t, "RunWithOutput")

	pather := NewPather(
		WithExecutor(executor),
		WithLookupEnv(lookupEnvOutWslMock),
	)

	got, err := pather.GetRealPath(path)
	r.Nil(err)
	assert.Equal(t, path, got)
}
