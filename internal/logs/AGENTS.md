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

## Cores ANSI Disponíveis

| Cor | Código | Uso |
|-----|--------|-----|
| Cyan/Bold | `\x1b[36m\x1b[1m` | Info (destaque) |
| Gray/Dim | `\x1b[90m` | Verbose (sutil) |
| Yellow | `\x1b[33m` | Debug |
| Reset | `\x1b[0m` | Fim da formatação |
