# VSCode Package

## Responsabilidade

Geração de URIs para VS Code Remote e integração com VS Code CLI.

## Padrões

- Funções puras, sem side-effects
- WSL path conversion quando `WSL_DISTRO_NAME` está setado
- URI format: `vscode-remote://dev-container+<hex_path><container_path>`
- Usar `paths.GetHostPath` para conversão de paths WSL

## Funções

| Função | Descrição |
|--------|-----------|
| `GetContainerURI(absPath, workspaceFolder)` | Gera URI para VS Code Remote |
| `BuildFolderURI(uri)` | Abre VS Code com folder-uri |

## URI Format

```
vscode-remote://dev-container+<hex_encoded_host_path><container_workspace_path>
```

Onde:
- `hex_encoded_host_path`: Path do host em hexadecimal
- `container_workspace_path`: Path dentro do container (ex: /workspaces)

## Uso

```go
import "github.com/Brennon-Oliveira/dev-cli/internal/vscode"

uri := vscode.GetContainerURI(absPath, workspaceFolder)
```
