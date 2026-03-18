# Completion Package

## Responsabilidade

Instalação de autocompletion para shells.

## Padrões

- Auto-detecção via `$SHELL` env var
- Append ao rc file se não existir
- Suporta: bash, zsh, powershell
- Usar constantes para paths de configuração

## Interface

```go
type Installer interface {
    Install(shell string) error
    DetectShell() string
}
```

## Shells Suportados

| Shell | RC File | Comando |
|-------|---------|---------|
| bash | ~/.bashrc | source completion.bash |
| zsh | ~/.zshrc | source completion.zsh |
| powershell | $PROFILE | . completion.ps1 |

## Adicionando Novo Shell

1. Adicionar detecção em `detectShell`
2. Criar função `install<Shell>`
3. Adicionar caso em `Install`
4. Atualizar documentação
