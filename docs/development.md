# Desenvolvimento e Contribuição

## Pré-requisitos

- Go 1.25.6 ou superior
- Docker ou Podman instalado
- Dev Container CLI (`npm install -g @devcontainers/cli`)
- VS Code instalado (se quiser testar os comandos `run` e `open`)

## Configuração Local

O CLI armazena configurações em `~/.dev-cli/config.json`. Você pode configurar a ferramenta de container manualmente:

```json
{
  "core": {
    "tool": "docker"
  }
}
```

## Compilando o Projeto

Para compilar o binário localmente:

```bash
go build -o dev .
```

Você pode então mover o binário para o seu PATH ou executá-lo diretamente: `./dev --help`.

## Testando

Certifique-se de validar as alterações em ambientes diferentes se estiver mexendo na lógica de caminhos ou processos, especialmente no **WSL**, que possui regras de mapeamento de rede e arquivos específicas.
