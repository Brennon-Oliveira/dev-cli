# Paths Package

## Responsabilidade

Resolução e conversão de caminhos de arquivos.

## Padrões

- `GetAbsPath` normaliza caminho para absoluto (usa `filepath.Abs`)
- Funções puras quando possível (sem side-effects)
- WSL detection via variável de ambiente `WSL_DISTRO_NAME`
- Conversão de paths WSL usando `wslpath` command

## Funções

| Função | Descrição |
|--------|-----------|
| `GetAbsPath(target)` | Normaliza caminho para absoluto |
| `GetHostPath(absPath)` | Resolve caminho considerando WSL |

## WSL Integration

Quando rodando em WSL (`WSL_DISTRO_NAME` setado), `GetHostPath` converte paths Linux para formato Windows usando `wslpath -w`.

## Uso

```go
import "github.com/Brennon-Oliveira/dev-cli/internal/paths"

absPath, err := paths.GetAbsPath(".")
hostPath := paths.GetHostPath(absPath)
```
