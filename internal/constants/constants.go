package constants

const (
	ToolDocker = "docker"
	ToolPodman = "podman"
)

const (
	LabelDevContainerFolder = "devcontainer.local_folder"
	LabelComposeProject     = "com.docker.compose.project"
)

const (
	DefaultWorkspaceFolder = "/workspaces"
	ConfigDirName          = ".dev-cli"
	ConfigFileName         = "config.json"
)

var ValidTools = []string{ToolDocker, ToolPodman}

var ValidBoolValues = []string{"true", "false"}
