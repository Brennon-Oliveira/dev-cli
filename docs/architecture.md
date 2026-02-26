# Arquitetura do Dev CLI

O Dev CLI é uma ferramenta escrita em Go projetada para orquestrar o ciclo de vida de Dev Containers, integrando-se nativamente com o VS Code e lidando com as complexidades de ambientes como o WSL.

## Visão Geral

A aplicação segue uma estrutura modular padrão em Go:

- `main.go`: Ponto de entrada que inicializa a execução dos comandos.
- `cmd/`: Contém a definição da interface de linha de comando usando a biblioteca **Cobra**. Cada arquivo aqui representa um comando ou subcomando.
- `internal/`: Contém a lógica de negócio protegida, inacessível por outros módulos externos.
    - `internal/config/`: Gerencia a configuração global da ferramenta (ex: qual engine de container usar, docker ou podman).
    - `internal/container/`: Implementa a lógica principal de interação com o Motor de Containers e com a CLI do `devcontainer`.

## Fluxo de Execução

1. O usuário executa um comando (ex: `dev run .`).
2. O `main.go` chama `cmd.Execute()`.
3. O Cobra identifica o comando correspondente em `cmd/run.go`.
4. O comando valida os argumentos e invoca funções em `internal/container`.
5. `internal/container` realiza chamadas de sistema para `devcontainer` ou `docker/podman`.
6. Se necessário, o VS Code é aberto usando uma URI customizada (`vscode-remote://...`).

## Integração WSL

Um dos grandes diferenciais desta ferramenta é o suporte robusto ao WSL. Ela detecta automaticamente se está rodando em ambiente Linux dentro do Windows e realiza a conversão de caminhos (usando `wslpath`) para garantir que o VS Code (rodando no Windows) consiga localizar o workspace corretamente.
