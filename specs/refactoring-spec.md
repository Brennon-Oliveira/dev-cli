# Especificação de Refatoração - Dev CLI

## Objetivo

Reestruturar o projeto Dev CLI para melhorar manutenibilidade, testabilidade e extensibilidade através de:
- Separação de domínios em pacotes distintos
- Introdução de interfaces (ports/adapters pattern)
- Padronização de tratamento de erros
- Centralização de constantes
- Sistema de logs estruturado
- Cobertura de testes com mocks

---

## Estrutura de Diretórios

### Atual
```
internal/
├── config/
│   ├── config.go
│   └── config_test.go
└── container/
    ├── paths.go
    ├── process.go
    ├── process_posix.go
    ├── process_windows.go
    └── container_test.go
```

### Nova
```
internal/
├── config/
│   ├── config.go
│   ├── handler.go
│   ├── config_test.go
│   └── AGENTS.md
├── container/
│   ├── client.go
│   ├── docker.go
│   ├── operations.go
│   ├── compose.go
│   ├── identifiers.go
│   ├── client_test.go
│   └── AGENTS.md
├── devcontainer/
│   ├── cli.go
│   ├── commands.go
│   ├── workspace.go
│   ├── cli_test.go
│   └── AGENTS.md
├── vscode/
│   ├── uri.go
│   ├── uri_test.go
│   └── AGENTS.md
├── exec/
│   ├── executor.go
│   ├── process.go
│   ├── process_posix.go
│   ├── process_windows.go
│   ├── exec_test.go
│   └── AGENTS.md
├── update/
│   ├── updater.go
│   ├── download.go
│   ├── extract.go
│   ├── update_test.go
│   └── AGENTS.md
├── completion/
│   ├── installer.go
│   ├── shells.go
│   ├── completion_test.go
│   └── AGENTS.md
├── paths/
│   ├── paths.go
│   ├── paths_test.go
│   └── AGENTS.md
├── logs/
│   ├── logs.go
│   ├── logs_debug.go
│   ├── logs_nodebug.go
│   ├── logs_test.go
│   └── AGENTS.md
└── constants/
    └── constants.go
```

---

## Interfaces

### 1. `internal/exec/executor.go`

```go
type Executor interface {
    Run(name string, args ...string) error
    RunInteractive(name string, args ...string) error
    RunDetached(name string, args ...string) error
    Output(name string, args ...string) (string, error)
}

type RealExecutor struct{}

func NewExecutor() *RealExecutor
func (e *RealExecutor) Run(name string, args ...string) error
func (e *RealExecutor) RunInteractive(name string, args ...string) error
func (e *RealExecutor) RunDetached(name string, args ...string) error
func (e *RealExecutor) Output(name string, args ...string) (string, error)
```

### 2. `internal/container/client.go`

```go
type ContainerClient interface {
    ListContainers(format string) error
    GetContainerID(filter string) (string, error)
    GetAllRelatedContainers(filter string) ([]string, error)
    StopContainers(ids []string) error
    RemoveContainers(ids []string) error
    ShowLogs(id string, follow bool) error
    ListPorts(id string) error
    CleanResources() error
    InspectLabel(id, label string) (string, error)
}

type DockerClient struct {
    tool     string
    executor exec.Executor
}

func NewClient(tool string, executor exec.Executor) *DockerClient
// ... implementa todos os métodos da interface
```

### 3. `internal/devcontainer/cli.go`

```go
type DevContainerCLI interface {
    Up(workspaceFolder string) error
    Exec(workspaceFolder string, command []string) error
    ReadConfiguration(workspaceFolder string) (*WorkspaceConfig, error)
}

type DevContainerCLIImpl struct {
    executor exec.Executor
}

func NewDevContainerCLI(executor exec.Executor) *DevContainerCLIImpl
// ... implementa todos os métodos da interface
```

### 4. `internal/update/updater.go`

```go
type Updater interface {
    CheckForUpdate() (*ReleaseInfo, error)
    Download(version string) (string, error)
    Extract(archivePath string) (string, error)
}

type GitHubUpdater struct {
    executor exec.Executor
}

type ReleaseInfo struct {
    TagName string
    URL     string
}
```

### 5. `internal/completion/installer.go`

```go
type Installer interface {
    Install(shell string) error
    DetectShell() string
}

type CompletionInstaller struct {
    rootCmd *cobra.Command
}
```

---

## Sistema de Logs

### Visão Geral

Sistema de logging estruturado com 3 níveis, projetado para ser simples e acessível de qualquer arquivo.

### API

```go
import "github.com/Brennon-Oliveira/dev-cli/internal/logs"

// Sempre aparece - informativo principal
logs.Info("Subindo container de desenvolvimento")

// Aparece com flag -v ou --verbose
logs.Verbose("Executando: docker ps --filter label=devcontainer.local_folder")

// Apenas em builds de desenvolvimento (build tag)
logs.Debug("Valor retornado: %v", someValue)
```

### Tipos de Log

| Tipo | Função | Quando Aparece | Formato |
|------|--------|----------------|---------|
| **Info** | `logs.Info(msg, args...)` | Sempre | Cor destacada (cyan/bold), prefixo com ícone |
| **Verbose** | `logs.Verbose(msg, args...)` | Flag `-v` ou `--verbose` | Cor sutil (gray/dim), prefixo simples |
| **Debug** | `logs.Debug(msg, args...)` | Apenas em build com tag `dev` | Cor amarela, prefixo `[DEBUG]` |

### Arquitetura

**`internal/logs/logs.go`:**
```go
package logs

import (
    "fmt"
    "os"
)

var verbose bool

func SetVerbose(v bool) {
    verbose = v
}

func Info(format string, args ...any) {
    msg := fmt.Sprintf(format, args...)
    fmt.Printf("\x1b[36m\x1b[1m➜\x1b[0m %s\n", msg)
}

func Verbose(format string, args ...any) {
    if !verbose {
        return
    }
    msg := fmt.Sprintf(format, args...)
    fmt.Printf("\x1b[90m  │ %s\x1b[0m\n", msg)
}
```

**`internal/logs/logs_debug.go` (build tag `dev`):**
```go
//go:build dev

package logs

import "fmt"

func Debug(format string, args ...any) {
    msg := fmt.Sprintf(format, args...)
    fmt.Printf("\x1b[33m[DEBUG] %s\x1b[0m\n", msg)
}
```

**`internal/logs/logs_nodebug.go` (build tag `!dev`):**
```go
//go:build !dev

package logs

func Debug(format string, args ...any) {
    // No-op em builds de produção
}
```

### Integração com Cobra

Adicionar flag global em `cmd/root.go`:
```go
var verboseFlag bool

func init() {
    rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Exibe logs detalhados")
}

func Execute() {
    logs.SetVerbose(verboseFlag)
    // ...
}
```

### Build Commands

```bash
# Build de produção (sem debug logs)
go build -o dev .

# Build de desenvolvimento (com debug logs)
go build -tags dev -o dev .
```

### Migração de fmt.Println

Substituir todas as ocorrências de `fmt.Println` e `fmt.Printf` por logs apropriados:

| Contexto | Usar |
|----------|------|
| "Subindo container em: %s" | `logs.Info()` |
| "Removendo containers parados..." | `logs.Info()` |
| Comando sendo executado | `logs.Verbose()` |
| Valor de variável para debug | `logs.Debug()` |
| Resultados/saídas para usuário | Manter `fmt.Print*` (não é log) |

---

## Migração por Domínio

### Fase 1: Constantes e Paths

| Arquivo Origem | Arquivo Destino | O que mover |
|----------------|-----------------|-------------|
| `container/process.go` (strings) | `constants/constants.go` | "docker", "podman", labels |
| `container/paths.go` | `paths/paths.go` | `GetAbsPath` |
| `container/paths.go` | `vscode/uri.go` | `GetHostPath`, `GetContainerURI` |
| `container/paths.go` | `devcontainer/workspace.go` | `getWorkspaceFolder` |

### Fase 2: Exec

| Arquivo Origem | Arquivo Destino | O que mover |
|----------------|-----------------|-------------|
| `container/process.go` | `exec/executor.go` | Interface Executor |
| `container/process.go` | `exec/process.go` | `ExecDetached`, `RunInteractive` (genérico) |
| `container/process_posix.go` | `exec/process_posix.go` | `applyDetachedAttr` |
| `container/process_windows.go` | `exec/process_windows.go` | `applyDetachedAttr` |

### Fase 3: DevContainer

| Arquivo Origem | Arquivo Destino | O que mover |
|----------------|-----------------|-------------|
| `container/process.go` | `devcontainer/commands.go` | `RunUpSync`, `RunInteractive` (devcontainer) |

### Fase 4: Container

| Arquivo Origem | Arquivo Destino | O que mover |
|----------------|-----------------|-------------|
| `container/process.go` | `container/client.go` | Interface ContainerClient |
| `container/process.go` | `container/docker.go` | DockerClient implementation |
| `container/process.go` | `container/operations.go` | `ListContainers`, `ShowLogs`, `ListPorts`, `CleanResources` |
| `container/process.go` | `container/compose.go` | `getAllRelatedContainers` |
| `container/process.go` | `container/identifiers.go` | `getContainerIDs`, `KillContainer`, `DownContainer` |

### Fase 5: Config Handler

| Arquivo Origem | Arquivo Destino | O que mover |
|----------------|-----------------|-------------|
| `cmd/config.go` | `config/handler.go` | `configHandler` struct, `handlers` map |

### Fase 6: Update

| Arquivo Origem | Arquivo Destino | O que mover |
|----------------|-----------------|-------------|
| `cmd/update.go` | `update/updater.go` | Lógica de update, interface |
| `cmd/update.go` | `update/download.go` | Download de releases |
| `cmd/update.go` | `update/extract.go` | `extractZip`, `extractTarGz` |

### Fase 7: Completion

| Arquivo Origem | Arquivo Destino | O que mover |
|----------------|-----------------|-------------|
| `cmd/add_completion.go` | `completion/installer.go` | Interface Installer |
| `cmd/add_completion.go` | `completion/shells.go` | `installZsh`, `installBash`, `installPowerShell`, `detectShell`, `appendToFileIfMissing` |

### Fase 8: CMD Refatoração

Após criar todos os pacotes internos, refatorar cada comando em `cmd/` para:
1. Usar interfaces injetadas
2. Remover lógica de negócio
3. Apenas orquestrar chamadas

### Fase 9: Migração de Logs

Substituir todas as ocorrências de `fmt.Println`/`fmt.Printf` por logs apropriados:

| Arquivo | fmt.Println Original | Substituir por |
|---------|---------------------|----------------|
| `container/process.go` | "Removendo containers parados..." | `logs.Info()` |
| `container/process.go` | "Removendo redes não utilizadas..." | `logs.Info()` |
| `container/process.go` | "Limpeza concluída." | `logs.Info()` |
| `container/process.go` | "Forçando parada..." | `logs.Info()` |
| `container/process.go` | "Parando graciosamente..." | `logs.Info()` |
| `cmd/up.go` | "Subindo container em: %s" | `logs.Info()` |
| `cmd/open.go` | "Abrindo VS Code..." | `logs.Info()` |
| `cmd/update.go` | "Verificando atualizações..." | `logs.Info()` |
| `cmd/update.go` | "Nova versão encontrada..." | `logs.Info()` |
| `cmd/update.go` | "Extraindo arquivos..." | `logs.Info()` |
| `cmd/add_completion.go` | "Autocompletar configurado..." | `logs.Info()` |
| Todos os arquivos | Comandos executados | `logs.Verbose()` |

---

## Correções Específicas

### 1. Remover Side-Effect em GetContainerURI

**Antes (`container/paths.go`):**
```go
func GetContainerURI(absPath string) string {
    // ...
    fmt.Println(final)  // SIDE EFFECT
    return final
}
```

**Depois (`vscode/uri.go`):**
```go
func GetContainerURI(absPath string) string {
    // ...
    return final  // SEM side-effect
}

// Se precisar imprimir, usar no cmd:
func PrintAndGetContainerURI(absPath string) string {
    uri := GetContainerURI(absPath)
    fmt.Println(uri)
    return uri
}
```

### 2. Tratamento de Erros em CleanResources

**Antes:**
```go
func CleanResources() error {
    fmt.Println("Removendo containers parados...")
    exec.Command(tool, "container", "prune", "-f").Run()  // Ignora erro
    
    fmt.Println("Removendo redes não utilizadas...")
    exec.Command(tool, "network", "prune", "-f").Run()    // Ignora erro
    
    fmt.Println("Limpeza concluída.")
    return nil
}
```

**Depois:**
```go
func (c *DockerClient) CleanResources() error {
    fmt.Println("Removendo containers parados...")
    if err := c.executor.Run(c.tool, "container", "prune", "-f"); err != nil {
        return fmt.Errorf("falha ao remover containers parados: %w", err)
    }
    
    fmt.Println("Removendo redes não utilizadas...")
    if err := c.executor.Run(c.tool, "network", "prune", "-f"); err != nil {
        return fmt.Errorf("falha ao remover redes: %w", err)
    }
    
    fmt.Println("Limpeza concluída.")
    return nil
}
```

### 3. Constantes Centralizadas

**Criar `internal/constants/constants.go`:**
```go
package constants

const (
    ToolDocker = "docker"
    ToolPodman = "podman"
    
    LabelDevContainerFolder = "devcontainer.local_folder"
    LabelComposeProject     = "com.docker.compose.project"
    
    DefaultWorkspaceFolder  = "/workspaces"
    
    ConfigDirName = ".dev-cli"
    ConfigFileName = "config.json"
)

var ValidTools = []string{ToolDocker, ToolPodman}
```

---

## Arquivos AGENTS.md

Criar `AGENTS.md` nos seguintes pacotes:

### `internal/container/AGENTS.md`
```markdown
# Container Package

## Responsabilidade
Gerenciar interações com Docker/Podman CLI.

## Padrões
- Sempre usar `ContainerClient` interface, nunca implementação direta
- Tool (docker/podman) vem do config via dependency injection
- Labels definidas em `internal/constants`

## Adicionando Nova Operação
1. Adicionar método na interface `ContainerClient`
2. Implementar em `DockerClient`
3. Criar teste com mock de `Executor`
```

### `internal/devcontainer/AGENTS.md`
```markdown
# DevContainer Package

## Responsabilidade
Wrappers para comandos da CLI `devcontainer`.

## Padrões
- Sempre usar `DevContainerCLI` interface
- Workspace folder resolvido via `ReadConfiguration` ou fallback

## Adicionando Novo Comando
1. Adicionar método na interface
2. Implementar usando `Executor`
3. Tratar erros com contexto
```

### `internal/exec/AGENTS.md`
```markdown
# Exec Package

## Responsabilidade
Abstração para execução de comandos externos.

## Padrões
- Interface `Executor` para todos os comandos externos
- Platform-specific code em arquivos separados (`_posix.go`, `_windows.go`)
- Nunca usar `exec.Command` diretamente fora deste pacote

## Testes
- Usar `MockExecutor` que registra comandos chamados
```

### `internal/config/AGENTS.md`
```markdown
# Config Package

## Responsabilidade
Gerenciamento de configuração persistente.

## Padrões
- Config salva em `~/.dev-cli/config.json`
- `Load()` retorna config com defaults
- `handler.go` gerencia validação e tipos de configuração

## Adicionando Nova Configuração
1. Adicionar campo na struct `GlobalConfig`
2. Criar handler em `handlers` map com validação
3. Atualizar testes de roundtrip
```

### `internal/vscode/AGENTS.md`
```markdown
# VSCode Package

## Responsabilidade
Geração de URIs para VS Code Remote.

## Padrões
- Funções puras, sem side-effects
- WSL path conversion quando `WSL_DISTRO_NAME` está setado
- URI format: `vscode-remote://dev-container+<hex_path><container_path>`
```

### `internal/paths/AGENTS.md`
```markdown
# Paths Package

## Responsabilidade
Resolução e conversão de caminhos.

## Padrões
- `GetAbsPath` normaliza caminho para absoluto
- WSL detection e conversão via `wslpath`
```

### `internal/update/AGENTS.md`
```markdown
# Update Package

## Responsabilidade
Verificação e download de atualizações.

## Padrões
- GitHub API para releases
- Download para temp dir
- Extração de tar.gz (unix) ou zip (windows)
- Usuário executa move final (sudo pode ser necessário)
```

### `internal/completion/AGENTS.md`
```markdown
# Completion Package

## Responsabilidade
Instalação de autocompletion para shells.

## Padrões
- Auto-detecção via `$SHELL` env var
- Append ao rc file se não existir
- Suporta: bash, zsh, powershell
```

### `internal/logs/AGENTS.md`
```markdown
# Logs Package

## Responsabilidade
Sistema de logging estruturado com múltiplos níveis.

## API
```go
logs.Info("mensagem")     // Sempre aparece
logs.Verbose("mensagem")  // Aparece com -v/--verbose
logs.Debug("mensagem")    // Apenas em build dev
```

## Padrões
- Nunca usar `fmt.Print*` para logs informativos
- `Info` para ações principais visíveis ao usuário
- `Verbose` para detalhes internos e comandos executados
- `Debug` apenas para desenvolvimento (build tag `dev`)

## Build
- Produção: `go build .` (Debug é no-op)
- Desenvolvimento: `go build -tags dev .` (Debug ativo)

## Adicionando Novo Formato de Log
1. Adicionar função em `logs.go`
2. Escolher cor ANSI apropriada
3. Documentar uso neste arquivo
```

---

## Cobertura de Testes

### Requisitos Mínimos

| Pacote | Cobertura Mínima | Tipo de Teste |
|--------|------------------|---------------|
| `config` | 80% | Unit + roundtrip |
| `container` | 75% | Unit com mock executor |
| `devcontainer` | 75% | Unit com mock executor |
| `vscode` | 80% | Unit |
| `exec` | 70% | Unit + integration |
| `update` | 75% | Unit com mock HTTP |
| `completion` | 70% | Unit |
| `paths` | 90% | Unit |
| `logs` | 85% | Unit (output capture) |
| `cmd` | 60% | Integration (executa comandos) |

### Testes com Mocks

**Mock Executor (`internal/exec/exec_test.go`):**
```go
type MockExecutor struct {
    Calls []string
    OutputErr error
    OutputResult string
}

func (m *MockExecutor) Run(name string, args ...string) error {
    m.Calls = append(m.Calls, fmt.Sprintf("%s %s", name, strings.Join(args, " ")))
    return nil
}
// ... outros métodos
```

**Testes de Logs (`internal/logs/logs_test.go`):**
```go
func TestInfo(t *testing.T) {
    var buf bytes.Buffer
    // Capturar stdout
    // Verificar se output contém mensagem formatada com cores
}

func TestVerbose_WhenDisabled(t *testing.T) {
    logs.SetVerbose(false)
    // Verificar que nada é impresso
}

func TestVerbose_WhenEnabled(t *testing.T) {
    logs.SetVerbose(true)
    // Verificar que mensagem é impressa com cor correta
}
```

### Testes de Integração

Manter testes existentes em `cmd/commands_test.go` que validam fluxo completo com executáveis mockados.

---

## Ordem de Implementação

1. **Criar constantes** - `internal/constants/constants.go`
2. **Criar logs** - `internal/logs/` (deve vir cedo pois será usado por outros pacotes)
3. **Criar paths** - Mover `GetAbsPath` de container para paths
4. **Criar exec interfaces** - `internal/exec/`
5. **Criar vscode** - Mover URI functions
6. **Criar devcontainer** - Separar devcontainer CLI
7. **Criar container interfaces** - Refatorar container package
8. **Criar config handler** - Separar lógica de handler
9. **Criar update** - Separar lógica de update
10. **Criar completion** - Separar lógica de completion
11. **Refatorar cmd/** - Usar novos pacotes + integrar flag verbose
12. **Mover/atualizar testes**
13. **Criar AGENTS.md em cada pacote**
14. **Atualizar documentação**

---

## Validação Final

- [ ] `go build ./...` compila sem erros
- [ ] `go build -tags dev ./...` compila sem erros
- [ ] `go test ./...` passa todos os testes
- [ ] `go vet ./...` sem warnings
- [ ] Cobertura de testes >= mínimos definidos
- [ ] Nenhum side-effect em funções que não deveriam ter
- [ ] Todos os erros tratados com contexto
- [ ] Constantes centralizadas
- [ ] Interfaces definidas para todos os serviços externos
- [ ] Sistema de logs implementado e integrado
- [ ] Flag `-v/--verbose` funcional
- [ ] Debug logs compilados condicionalmente via build tag
