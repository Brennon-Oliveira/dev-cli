# Dev CLI

[English](README.md) | [Português](README.pt.md)

A command-line interface for managing Dev Container lifecycle and native VS Code integration in detached mode.

## 📥 Installation

### Linux and macOS (Native WSL Support)

#### Option 1: Snap Installation (Recommended)

```bash
sudo snap install dev-cli --classic
dev-cli --help
```

#### Option 2: Direct Download

Replace `[OS]` with `linux` or `macos` and `[ARCH]` with `amd64` or `arm64` according to your architecture.

```bash
# 1. Download the compressed artifact
curl -LO https://github.com/Brennon-Oliveira/dev-cli/releases/latest/download/dev-[OS]-[ARCH].tar.gz

# 2. Extract the file
tar -xzf dev-[OS]-[ARCH].tar.gz

# 3. Move the executable to your system PATH
sudo mv dev /usr/local/bin/dev-cli

# 4. Clean up the downloaded file
rm dev-[OS]-[ARCH].tar.gz

# 5. Verify the installation
dev-cli --help
```

### Windows

1. Download `dev-windows-amd64.zip` from [latest release](https://github.com/Brennon-Oliveira/dev-cli/releases/latest)
2. Extract the contents of the `.zip` file
3. Rename `dev.exe` to `dev-cli.exe`
4. Move the `dev-cli.exe` binary to a secure directory (e.g., `C:\Tools\bin`)
5. Add this directory to your Windows `PATH` environment variable
6. Verify installation: `dev-cli --help`

## 🛠️ Commands

### Container Lifecycle

- **`dev-cli run [path]`** (Recommended) - Provisions the container and immediately opens VS Code in the mapped directory
- **`dev-cli up [path]`** - Provisions and starts the dev container in the background without opening the editor
- **`dev-cli open [path]`** - Opens VS Code directly connected to an already running dev container, dynamically resolving the `workspaceFolder` from `devcontainer.json`
- **`dev-cli kill [path]`** - Instantly locates and terminates the container process attached to the target workspace
- **`dev-cli down [path]`** - Gracefully stops the container of the current workspace

### Environment Interaction

- **`dev-cli shell [path]`** - Injects an interactive shell (`zsh`, `bash`, or `sh`) directly into the active container
- **`dev-cli exec "command"`** - Passes commands and parameters for execution in the isolated container context (e.g., `dev-cli exec npm run build`)

### Monitoring and Diagnostics

- **`dev-cli list`** or **`dev-cli info`** - Returns a list of all dev containers running on the local host
- **`dev-cli logs [path]`** - Displays the container's standard output. Use the `-f` flag for real-time monitoring (*tail*)
- **`dev-cli ports [path]`** - Lists all active network mappings and exposed ports between the host and the current container

### Configuration

- **`dev-cli config [key] [value]`** - Manages CLI configuration settings (e.g., Docker vs Podman selection)
- **`dev-cli add-completion [bash|zsh|powershell]`** - Automatically configures shell auto-completion

### Maintenance

- **`dev-cli clean`** - Performs Docker resource cleanup by removing stopped containers and orphaned networks
- **`dev-cli update`** - (Experimental) Downloads the latest CLI version and prepares for installation

## ⚙️ Use Cases

### Immediate Onboarding
After cloning a project, there's no need to open VS Code, locate the folder, and click "Reopen in Container". Simply run `dev-cli run` at the repository root from the terminal. The CLI resolves the Docker build, handles path interoperability (if using WSL), and injects `code` into the final structure.

```bash
cd my-project
dev-cli run
```

### Headless Execution
If you only need to run tests or compile artifacts in a standardized environment, use `dev-cli up` to bring the infrastructure up invisibly and `dev-cli exec` to trigger routines, consuming less system memory by not instantiating Electron.

```bash
dev-cli up .
dev-cli exec npm run test
dev-cli exec npm run build
```

### Advanced Path Resolution
The project natively analyzes `devcontainer.json` configurations via regex, ensuring the editor accesses the correct workspace root (`workspaceFolder`), automatically handling fallbacks, short paths, and VS Code URI parse bugs.

### Interactive Terminal Sessions
Quickly access container shell for debugging or manual operations:

```bash
dev-cli shell
npm install
npm run dev
```

### Container Management
Monitor and control multiple development containers:

```bash
dev-cli list                    # See all running containers
dev-cli logs . -f               # Follow logs in real-time
dev-cli ports .                 # Check port mappings
dev-cli kill .                  # Stop and remove container
```

## 🔧 Configuration

### Container Engine Selection

By default, Dev CLI uses Docker. To use Podman:

```bash
dev-cli config --global core.tool podman
```

View current configuration:

```bash
dev-cli config --global core.tool
```

### Shell Completion

Install shell auto-completion for your shell:

```bash
# Bash
dev-cli add-completion bash

# Zsh
dev-cli add-completion zsh

# PowerShell
dev-cli add-completion powershell
```

## 📋 Requirements

- **Docker** or **Podman** installed
- **Dev Container CLI** (auto-detected if available): `npm install -g @devcontainers/cli`
- **VS Code** (optional, for `run` and `open` commands)
- Go 1.25.6+ (for building from source)

## 🚀 Quick Start

```bash
# 1. Clone a repository
git clone https://github.com/example/project.git
cd project

# 2. Bring up the dev container and open in VS Code
dev-cli run

# 3. Or, start the container in the background
dev-cli up
dev-cli shell
npm install
npm run dev
```

## 🔍 Troubleshooting

### "Container not found" Error

Ensure you're in the project directory that contains `devcontainer.json`:

```bash
dev-cli run .
```

### Port Already in Use

List current port mappings and check for conflicts:

```bash
dev-cli ports
```

### WSL Path Issues

The CLI automatically handles WSL path conversion. If you experience issues:

1. Verify you're running in WSL: `uname -a`
2. Ensure Docker Desktop is running and configured for WSL
3. Run with verbose output: `dev-cli run --verbose .`

### Container Won't Start

Check container logs:

```bash
dev-cli logs . -f
```

Check Docker daemon:

```bash
docker ps -a
```

## 🤝 Contributing

We welcome contributions! See our development guide in `docs/development.md` for:

- Building from source
- Running tests
- Code patterns and conventions
- WSL-specific considerations

## 📚 Documentation

- **[Architecture](docs/architecture.md)** - System design and components
- **[Commands](docs/commands.md)** - Command creation guide
- **[Patterns](docs/patterns.md)** - Development patterns
- **[Development](docs/development.md)** - Setup and testing
- **[cmd/AGENTS.md](cmd/AGENTS.md)** - Command structure guide
- **[internal/AGENTS.md](internal/AGENTS.md)** - Internal package structure

## 📝 License

This project is open source and available under the MIT License.

## 🔗 Links

- [GitHub Repository](https://github.com/Brennon-Oliveira/dev-cli)
- [Releases](https://github.com/Brennon-Oliveira/dev-cli/releases)
- [Dev Container Specification](https://containers.dev)

## ✨ Features

- ✅ One-command container setup and VS Code integration
- ✅ Full WSL support with automatic path conversion
- ✅ Support for both Docker and Podman
- ✅ Real-time log monitoring
- ✅ Port mapping visibility
- ✅ Shell completion for bash, zsh, and PowerShell
- ✅ Configuration management
- ✅ Experimental self-update functionality
