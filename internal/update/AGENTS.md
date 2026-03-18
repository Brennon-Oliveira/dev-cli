# Update Package

## Responsabilidade

Verificação e download de atualizações do CLI.

## Padrões

- GitHub API para releases
- Download para temp dir
- Extração de tar.gz (unix) ou zip (windows)
- Usuário executa move final (sudo pode ser necessário)
- Usar `Executor` para comandos externos

## Interface

```go
type Updater interface {
    CheckForUpdate(currentVersion string) (*ReleaseInfo, error)
    Download(version string) (string, error)
    Extract(archivePath string) (string, error)
}
```

## Fluxo

1. `CheckForUpdate` consulta GitHub API
2. `Download` baixa o artefato para temp
3. `Extract` extrai o binário
4. Usuário move manualmente (pode precisar de sudo)

## Build Tags

- Produção: build normal
- Desenvolvimento: pode mockar HTTP client
