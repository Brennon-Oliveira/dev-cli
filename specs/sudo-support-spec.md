# Especificação - Suporte a Sudo para Docker/Podman

## Objetivo

Adicionar configuração `core.use-sudo` que permite executar todos os comandos docker/podman com `sudo`, resolvendo problemas de permissão sem exigir que o usuário esteja no grupo docker.

---

## Visão Geral

```
┌─────────────────────────────────────────────────────────────┐
│                      Configuração                           │
│  ~/.dev-cli/config.json                                     │
│  { "core": { "tool": "docker", "use-sudo": true } }         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      DockerClient                           │
│  buildArgs("ps", "-a")                                      │
│    → ["sudo", "docker", "ps", "-a"]  (se useSudo=true)      │
│    → ["docker", "ps", "-a"]           (se useSudo=false)    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Detecção de Erro                          │
│  "permission denied" → Erro customizado com sugestão        │
└─────────────────────────────────────────────────────────────┘
```

---

## Arquivos a Alterar

### 1. `internal/config/config.go`

**Adicionar campo na struct:**

```go
type GlobalConfig struct {
    Core struct {
        Tool    string `json:"tool"`
        UseSudo bool   `json:"useSudo"`
    } `json:"core"`
}
```

**Nota:** Não precisa de valor default (false é o zero value de bool)

---

### 2. `internal/config/handler.go`

**Adicionar handler para nova configuração:**

```go
var Handlers = map[string]ConfigHandler{
    "core.tool": {
        ValidValues: constants.ValidTools,
        Label:       "Selecione o motor de containers padrão",
        Get: func(cfg GlobalConfig) string {
            return cfg.Core.Tool
        },
        Set: func(cfg *GlobalConfig, val string) {
            cfg.Core.Tool = val
        },
    },
    "core.use-sudo": {
        ValidValues: constants.ValidBoolValues,
        Label:       "Usar sudo para comandos Docker/Podman?",
        Get: func(cfg GlobalConfig) string {
            return fmt.Sprintf("%v", cfg.Core.UseSudo)
        },
        Set: func(cfg *GlobalConfig, val string) {
            cfg.Core.UseSudo = val == "true"
        },
    },
}
```

---

### 3. `internal/constants/constants.go`

**Adicionar valores válidos para boolean:**

```go
var ValidBoolValues = []string{"true", "false"}
```

---

### 4. `internal/container/client.go`

**Atualizar struct e construtor:**

```go
type DockerClient struct {
    tool     string
    useSudo  bool
    executor exec.Executor
}

func NewDockerClient(tool string, useSudo bool, executor exec.Executor) *DockerClient {
    return &DockerClient{
        tool:     tool,
        useSudo:  useSudo,
        executor: executor,
    }
}
```

**Adicionar método helper:**

```go
func (d *DockerClient) buildArgs(args ...string) []string {
    if d.useSudo {
        return append([]string{"sudo", d.tool}, args...)
    }
    return append([]string{d.tool}, args...)
}
```

**Adicionar erro customizado:**

```go
var ErrPermissionDenied = errors.New("permissão negada ao acessar Docker/Podman")

func wrapPermissionError(err error) error {
    if err == nil {
        return nil
    }
    
    errStr := err.Error()
    if strings.Contains(errStr, "permission denied") ||
       strings.Contains(errStr, "Permission denied") ||
       strings.Contains(errStr, "Cannot connect to the Docker daemon") {
        return fmt.Errorf("%w\n\nDica: execute 'dev config core.use-sudo true --global' para usar sudo", ErrPermissionDenied)
    }
    
    return err
}
```

**Atualizar MockContainerClient:**

```go
type MockContainerClient struct {
    // ... campos existentes
    UseSudo bool  // adicionar se necessário para verificação
}
```

---

### 5. `internal/container/operations.go`

**Atualizar todos os métodos para usar buildArgs:**

```go
func (d *DockerClient) ListContainers() error {
    logs.Info("Listando containers de desenvolvimento")
    format := "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Label \"devcontainer.local_folder\"}}"
    logs.Verbose("executando: %s ps --filter label=%s --format ...", d.tool, constants.LabelDevContainerFolder)
    
    args := d.buildArgs("ps", "--filter", "label="+constants.LabelDevContainerFolder, "--format", format)
    return wrapPermissionError(d.executor.Run(args[0], args[1:]...))
}

func (d *DockerClient) CleanResources() error {
    logs.Info("Removendo containers parados")
    logs.Verbose("executando: %s container prune -f", d.tool)
    
    args := d.buildArgs("container", "prune", "-f")
    if err := wrapPermissionError(d.executor.Run(args[0], args[1:]...)); err != nil {
        return fmt.Errorf("falha ao remover containers parados: %w", err)
    }

    logs.Info("Removendo redes não utilizadas")
    logs.Verbose("executando: %s network prune -f", d.tool)
    
    args = d.buildArgs("network", "prune", "-f")
    if err := wrapPermissionError(d.executor.Run(args[0], args[1:]...)); err != nil {
        return fmt.Errorf("falha ao remover redes: %w", err)
    }

    logs.Success("Limpeza concluída")
    return nil
}

// Aplicar mesmo padrão para:
// - ShowLogs
// - ListPorts
// - StopContainers
// - RemoveContainers
```

---

### 6. `internal/container/compose.go`

**Atualizar para usar buildArgs:**

```go
func (d *DockerClient) GetAllRelatedContainers(path string) ([]string, error) {
    // ...
    for _, p := range pathsToTry {
        filter := fmt.Sprintf("label=%s=%s", constants.LabelDevContainerFolder, p)
        logs.Verbose("executando: %s ps -a -q --filter %s", d.tool, filter)

        args := d.buildArgs("ps", "-a", "-q", "--filter", filter)
        out, err := d.executor.Output(args[0], args[1:]...)
        if err != nil {
            if permErr := wrapPermissionError(err); permErr != err {
                return nil, permErr
            }
            continue
        }
        // ...
    }
    // ...
}
```

---

### 7. `internal/container/identifiers.go`

**Atualizar para usar buildArgs:**

```go
func (d *DockerClient) GetContainerID(path string) (string, error) {
    // ...
    for _, p := range pathsToTry {
        filter := fmt.Sprintf("label=%s=%s", constants.LabelDevContainerFolder, p)
        logs.Verbose("executando: %s ps -q --filter %s", d.tool, filter)

        args := d.buildArgs("ps", "-q", "--filter", filter)
        out, err := d.executor.Output(args[0], args[1:]...)
        if err != nil {
            if permErr := wrapPermissionError(err); permErr != err {
                return "", permErr
            }
            continue
        }
        // ...
    }
    // ...
}

func (d *DockerClient) InspectLabel(id, label string) (string, error) {
    format := fmt.Sprintf("{{ if .Config.Labels }}{{ index .Config.Labels %q }}{{ end }}", label)
    logs.Verbose("executando: %s inspect -f '%s' %s", d.tool, format, id)

    args := d.buildArgs("inspect", "-f", format, id)
    out, err := d.executor.Output(args[0], args[1:]...)
    if err != nil {
        return "", wrapPermissionError(err)
    }

    return strings.TrimSpace(out), nil
}
```

---

### 8. Arquivos `cmd/*.go`

**Arquivos a atualizar:**
- `clean.go`
- `down.go`
- `kill.go`
- `list.go`
- `logs.go`
- `ports.go`
- `shell.go`

**Padrão de alteração:**

```go
// Antes
client := container.NewDockerClient(cfg.Core.Tool, executor)

// Depois
client := container.NewDockerClient(cfg.Core.Tool, cfg.Core.UseSudo, executor)
```

---

## Fluxo de Uso

### Configuração

```bash
# Habilitar sudo
dev config core.use-sudo true --global

# Verificar configuração
dev config core.use-sudo --global

# Desabilitar sudo
dev config core.use-sudo false --global
```

### Execução

```bash
# Com use-sudo: true
dev list
# Executa: sudo docker ps --filter label=devcontainer.local_folder ...

# Com use-sudo: false (padrão)
dev list
# Executa: docker ps --filter label=devcontainer.local_folder ...
```

### Erro de Permissão

```bash
# Sem sudo configurado e sem permissão
$ dev list
Error: permissão negada ao acessar Docker/Podman

Dica: execute 'dev config core.use-sudo true --global' para usar sudo
```

---

## Testes

### `internal/config/config_test.go`

```go
func TestLoad_DefaultUseSudo(t *testing.T) {
    cfg := Load()
    if cfg.Core.UseSudo {
        t.Error("expected UseSudo to be false by default")
    }
}

func TestSaveAndLoad_UseSudo(t *testing.T) {
    // ... criar config com UseSudo: true, salvar, carregar e verificar
}
```

### `internal/config/handler_test.go` (novo)

```go
func TestHandler_UseSudo(t *testing.T) {
    handler, exists := Handlers["core.use-sudo"]
    if !exists {
        t.Fatal("handler not found")
    }
    
    if !reflect.DeepEqual(handler.ValidValues, []string{"true", "false"}) {
        t.Error("invalid valid values")
    }
    
    cfg := GlobalConfig{}
    handler.Set(&cfg, "true")
    if !cfg.Core.UseSudo {
        t.Error("expected UseSudo to be true")
    }
    
    if handler.Get(cfg) != "true" {
        t.Error("expected Get to return 'true'")
    }
}
```

### `internal/container/client_test.go`

```go
func TestDockerClient_BuildArgs_WithSudo(t *testing.T) {
    mock := exec.NewMockExecutor()
    client := NewDockerClient("docker", true, mock)
    
    client.ListContainers()
    
    if len(mock.Calls) != 1 {
        t.Fatal("expected 1 call")
    }
    
    if !strings.HasPrefix(mock.Calls[0], "sudo docker") {
        t.Errorf("expected sudo prefix, got: %s", mock.Calls[0])
    }
}

func TestDockerClient_BuildArgs_WithoutSudo(t *testing.T) {
    mock := exec.NewMockExecutor()
    client := NewDockerClient("docker", false, mock)
    
    client.ListContainers()
    
    if len(mock.Calls) != 1 {
        t.Fatal("expected 1 call")
    }
    
    if strings.HasPrefix(mock.Calls[0], "sudo") {
        t.Errorf("expected no sudo prefix, got: %s", mock.Calls[0])
    }
}
```

---

## Ordem de Implementação

1. **Constantes** - `internal/constants/constants.go`
   - Adicionar `ValidBoolValues`

2. **Config** - `internal/config/`
   - Adicionar campo `UseSudo` na struct
   - Adicionar handler para `core.use-sudo`
   - Atualizar testes

3. **Container** - `internal/container/`
   - Adicionar campo `useSudo` e atualizar construtor
   - Criar método `buildArgs()`
   - Criar função `wrapPermissionError()`
   - Atualizar todos os métodos que executam comandos
   - Atualizar testes

4. **CMD** - `cmd/*.go`
   - Atualizar todos os comandos que usam `NewDockerClient`

5. **Testes de Integração**
   - Verificar fluxo completo
   - Verificar mensagens de erro

---

## Validação Final

- [ ] `go build ./...` compila sem erros
- [ ] `go test ./...` passa todos os testes
- [ ] `dev config core.use-sudo true --global` funciona
- [ ] `dev config core.use-sudo --global` retorna valor correto
- [ ] Comandos docker executam com `sudo` quando configurado
- [ ] Comandos docker executam sem `sudo` por padrão
- [ ] Erro de permissão mostra mensagem com sugestão
- [ ] Logs verbose mostram comando completo com `sudo`
