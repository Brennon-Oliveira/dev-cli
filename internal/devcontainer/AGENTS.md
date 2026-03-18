# DevContainer Package

## Responsabilidade

Wrappers para comandos da CLI `devcontainer`.

## Padrões

- Sempre usar `DevContainerCLI` interface
- Workspace folder resolvido via `ReadConfiguration` ou fallback
- Usar `Executor` para comandos externos
- Logs via pacote `logs`

## Interface

```go
type DevContainerCLI interface {
    Up(workspaceFolder string) error
    Exec(workspaceFolder string, command []string) error
    ReadConfiguration(workspaceFolder string) (*WorkspaceConfig, error)
}
```

## Adicionando Novo Comando

1. Adicionar método na interface
2. Implementar usando `Executor`
3. Tratar erros com contexto
4. Adicionar logs apropriados
5. Criar teste com mock executor

## Comandos Disponíveis

| Método | CLI Command | Descrição |
|--------|-------------|-----------|
| `Up` | `devcontainer up` | Sobe o container |
| `Exec` | `devcontainer exec` | Executa comando no container |
| `ReadConfiguration` | `devcontainer read-configuration` | Lê configuração do workspace |
