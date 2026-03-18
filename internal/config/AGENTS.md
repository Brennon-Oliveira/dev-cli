# Config Package

## Responsabilidade

Gerenciamento de configuração persistente.

## Padrões

- Config salva em `~/.dev-cli/config.json`
- `Load()` retorna config com defaults
- `handler.go` gerencia validação e tipos de configuração
- Usar constantes de `internal/constants` para nomes de arquivos

## Adicionando Nova Configuração

1. Adicionar campo na struct `GlobalConfig`
2. Criar handler em `handlers` map com validação
3. Atualizar testes de roundtrip
4. Atualizar documentação

## Handler Structure

```go
type ConfigHandler struct {
    ValidValues []string
    Label       string
    Get         func(cfg GlobalConfig) string
    Set         func(cfg *GlobalConfig, val string)
}
```

## Configurações Disponíveis

| Chave | Valores | Descrição |
|-------|---------|-----------|
| `core.tool` | docker, podman | Motor de containers |
