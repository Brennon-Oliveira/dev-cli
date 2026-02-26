# Padrões e Convenções

Para manter a consistência e a manutenibilidade do projeto, siga estes padrões:

## Organização de Código

- **Encapsulamento**: Lógica que não precisa ser exposta deve ficar em `internal/`.
- **Engine de Containers**: Nunca chame `docker` ou `podman` diretamente com strings estáticas. Use `config.Load().Core.Tool` para respeitar a preferência do usuário.
- **Chamadas de Sistema**: Utilize as funções auxiliares em `internal/container/process.go` (como `ExecDetached` ou `RunInteractive`) para garantir que os processos sejam tratados corretamente em diferentes sistemas operacionais.

## Tratamento de Erros

- Retorne erros em vez de dar `log.Fatal` no meio do código (exceto em casos extremos no `main.go` ou `root.go`).
- Envolva erros com contexto se necessário, usando `fmt.Errorf("contexto: %w", err)`.

## Compatibilidade de SO

O projeto suporta Windows, Linux e macOS, além do WSL.
- Se precisar de comportamento específico por SO, use arquivos com sufixos (ex: `process_windows.go`, `process_posix.go`).
- Use `filepath.Join` em vez de concatenação manual de strings com `/` ou `\`.

## Funções Importantes em `internal/container`

- `GetAbsPath(target string)`: Normaliza o caminho para absoluto.
- `GetHostPath(absPath string)`: Resolve o caminho para o formato do Host (crucial para WSL).
- `GetContainerURI(absPath string)`: Gera a URI de conexão remota do VS Code.
- `getAllRelatedContainers(path string)`: Localiza todos os containers (incluindo serviços do Docker Compose) vinculados a um projeto.

## Estendendo Configurações

Para adicionar novas configurações ao CLI:
1.  Adicione o campo correspondente na struct `GlobalConfig` em `internal/config/config.go`.
2.  Defina um novo `handler` no mapa `handlers` em `cmd/config.go`.
3.  O `handler` deve definir os valores válidos, a mensagem de exibição e as funções de `Get` e `Set`.
