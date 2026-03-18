# Container Package

## Responsabilidade

Gerenciar interações com Docker/Podman CLI.

## Padrões

- Sempre usar `ContainerClient` interface, nunca implementação direta
- Tool (docker/podman) vem do config via dependency injection
- Labels definidas em `internal/constants`
- Usar `Executor` para comandos externos
- Logs via pacote `logs`

## Interface

```go
type ContainerClient interface {
    ListContainers() error
    GetContainerID(path string) (string, error)
    GetAllRelatedContainers(path string) ([]string, error)
    StopContainers(ids []string) error
    RemoveContainers(ids []string) error
    ShowLogs(path string, follow bool) error
    ListPorts(path string) error
    CleanResources() error
}
```

## Adicionando Nova Operação

1. Adicionar método na interface `ContainerClient`
2. Implementar em `DockerClient`
3. Criar teste com mock de `Executor`
4. Adicionar logs apropriados

## Estrutura de Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `client.go` | Interface e mocks |
| `docker.go` | Implementação DockerClient |
| `operations.go` | Operações de container (list, logs, ports, clean) |
| `compose.go` | Lógica de Docker Compose |
| `identifiers.go` | Identificação de containers por path |
