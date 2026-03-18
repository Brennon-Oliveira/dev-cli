# Exec Package

## Responsabilidade

Abstração para execução de comandos externos.

## Padrões

- Interface `Executor` para todos os comandos externos
- Platform-specific code em arquivos separados (`_posix.go`, `_windows.go`)
- Nunca usar `exec.Command` diretamente fora deste pacote
- Injetar Executor via construtor para permitir testes com mocks

## Interface

```go
type Executor interface {
    Run(name string, args ...string) error
    RunInteractive(name string, args ...string) error
    RunDetached(name string, args ...string) error
    Output(name string, args ...string) (string, error)
}
```

## Testes

- Usar `MockExecutor` que registra comandos chamados
- `MockExecutor` pode ser configurado para retornar erros específicos

## Adicionando Novo Método

1. Adicionar método na interface `Executor`
2. Implementar em `RealExecutor`
3. Adicionar no `MockExecutor` para testes
4. Atualizar testes existentes
