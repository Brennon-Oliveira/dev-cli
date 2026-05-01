# Dev CLI

Uma interface de linha de comando para gerenciar o ciclo de vida de Dev Containers e integração nativa com VS Code em modo *detached*.

## 📥 Instalação

### Linux e macOS (Compatível nativamente com WSL)

#### Opção 1: Instalação via Snap (Recomendado)

```bash
sudo snap install dev-cli --classic
dev-cli --help
```

#### Opção 2: Download Direto

Substitua `[OS]` por `linux` ou `macos` e `[ARCH]` por `amd64` ou `arm64` de acordo com a sua arquitetura.

```bash
# 1. Baixe o artefato compactado
curl -LO https://github.com/Brennon-Oliveira/dev-cli/releases/latest/download/dev-[OS]-[ARCH].tar.gz

# 2. Descompacte o arquivo
tar -xzf dev-[OS]-[ARCH].tar.gz

# 3. Mova o executável para o PATH do sistema
sudo mv dev /usr/local/bin/dev-cli

# 4. Limpe o arquivo baixado
rm dev-[OS]-[ARCH].tar.gz

# 5. Valide a instalação
dev-cli --help
```

### Windows

1. Baixe o arquivo `dev-windows-amd64.zip` na [última release](https://github.com/Brennon-Oliveira/dev-cli/releases/latest)
2. Extraia o conteúdo do arquivo `.zip`
3. Renomeie `dev.exe` para `dev-cli.exe`
4. Mova o binário `dev-cli.exe` para um diretório seguro (ex: `C:\Ferramentas\bin`)
5. Adicione este diretório à variável de ambiente `PATH` do Windows
6. Valide a instalação: `dev-cli --help`

## 🛠️ Comandos

### Ciclo de Vida do Container

- **`dev-cli run [caminho]`** (Recomendado) - Provisiona o container e imediatamente abre o VS Code no diretório mapeado
- **`dev-cli up [caminho]`** - Provisiona e inicia o dev container em segundo plano, sem abrir o editor
- **`dev-cli open [caminho]`** - Abre o VS Code diretamente conectado ao dev container já em execução, resolvendo dinamicamente o `workspaceFolder` do `devcontainer.json`
- **`dev-cli kill [caminho]`** - Localiza e encerra instantaneamente o processo do container atrelado ao workspace alvo
- **`dev-cli down [caminho]`** - Para graciosamente o container do workspace atual

### Interação com o Ambiente

- **`dev-cli shell [caminho]`** - Injeta um shell interativo (`zsh`, `bash` ou `sh`) diretamente dentro do container ativo
- **`dev-cli exec "comando"`** - Repassa comandos e parâmetros para execução no contexto isolado do container (ex: `dev-cli exec npm run build`)

### Monitoramento e Diagnóstico

- **`dev-cli list`** ou **`dev-cli info`** - Retorna a lista de todos os dev containers em execução no host local
- **`dev-cli logs [caminho]`** - Exibe a saída padrão do container. Use a flag `-f` para acompanhamento em tempo real (*tail*)
- **`dev-cli ports [caminho]`** - Lista todos os mapeamentos de rede e portas expostas ativas entre o host e o container atual

### Configuração

- **`dev-cli config [chave] [valor]`** - Gerencia as configurações da CLI (ex: seleção de Docker vs Podman)
- **`dev-cli add-completion [bash|zsh|powershell]`** - Configura o autocompletar da CLI automaticamente no seu shell

### Manutenção

- **`dev-cli clean`** - Realiza a liberação de recursos do Docker, removendo containers parados e redes órfãs
- **`dev-cli update`** - (EXPERIMENTAL) Baixa a última versão da CLI e prepara para instalação

## ⚙️ Casos de Uso

### Onboarding Imediato
Após clonar um projeto, não é necessário abrir o VS Code, localizar a pasta e clicar em "Reopen in Container". Basta rodar `dev-cli run` na raiz do repositório pelo terminal. A CLI resolve o build do Docker, lida com a interoperabilidade de caminhos (caso esteja usando WSL) e injeta o `code` na estrutura final.

```bash
cd meu-projeto
dev-cli run
```

### Execução Headless
Se você precisa apenas rodar testes ou compilar artefatos em um ambiente padronizado, utilize `dev-cli up` para subir a infraestrutura invisível e `dev-cli exec` para acionar as rotinas, consumindo menos memória do sistema host por não instanciar o Electron.

```bash
dev-cli up .
dev-cli exec npm run test
dev-cli exec npm run build
```

### Resolução Avançada de Caminhos
O projeto analisa as configurações do `devcontainer.json` nativamente através de regex, garantindo que o editor acesse a pasta raiz real (`workspaceFolder`), lidando automaticamente com fallbacks, caminhos curtos e bugs de parse de URI do VS Code.

### Sessões de Terminal Interativo
Acesse rapidamente o shell do container para depuração ou operações manuais:

```bash
dev-cli shell
npm install
npm run dev
```

### Gerenciamento de Containers
Monitore e controle múltiplos containers de desenvolvimento:

```bash
dev-cli list                    # Veja todos os containers em execução
dev-cli logs . -f               # Acompanhe os logs em tempo real
dev-cli ports .                 # Verifique os mapeamentos de portas
dev-cli kill .                  # Encerre e remova o container
```

## 🔧 Configuração

### Seleção do Motor de Container

Por padrão, Dev CLI usa Docker. Para usar Podman:

```bash
dev-cli config --global core.tool podman
```

Visualize a configuração atual:

```bash
dev-cli config --global core.tool
```

### Autocompletar do Shell

Instale o autocompletar para seu shell:

```bash
# Bash
dev-cli add-completion bash

# Zsh
dev-cli add-completion zsh

# PowerShell
dev-cli add-completion powershell
```

## 📋 Requisitos

- **Docker** ou **Podman** instalado
- **Dev Container CLI** (auto-detectado se disponível): `npm install -g @devcontainers/cli`
- **VS Code** (opcional, para comandos `run` e `open`)
- Go 1.25.6+ (para compilar a partir do código-fonte)

## 🚀 Início Rápido

```bash
# 1. Clone um repositório
git clone https://github.com/exemplo/projeto.git
cd projeto

# 2. Traga o dev container e abra no VS Code
dev-cli run

# 3. Ou, inicie o container em background
dev-cli up
dev-cli shell
npm install
npm run dev
```

## 🔍 Resolução de Problemas

### Erro "Container não encontrado"

Certifique-se de que você está no diretório do projeto que contém `devcontainer.json`:

```bash
dev-cli run .
```

### Porta Já em Uso

Liste os mapeamentos de porta atuais e verifique conflitos:

```bash
dev-cli ports
```

### Problemas de Caminho no WSL

A CLI manipula automaticamente a conversão de caminhos WSL. Se você experimentar problemas:

1. Verifique se você está rodando no WSL: `uname -a`
2. Garanta que Docker Desktop esteja rodando e configurado para WSL
3. Execute com saída verbosa: `dev-cli run --verbose .`

### Container Não Inicia

Verifique os logs do container:

```bash
dev-cli logs . -f
```

Verifique o daemon Docker:

```bash
docker ps -a
```

## 🤝 Contribuindo

Contribuições são bem-vindas! Veja nosso guia de desenvolvimento em `docs/development.md` para:

- Compilar a partir do código-fonte
- Executar testes
- Padrões de código e convenções
- Considerações específicas de WSL

## 📚 Documentação

- **[Arquitetura](docs/architecture.md)** - Design e componentes do sistema
- **[Comandos](docs/commands.md)** - Guia de criação de comandos
- **[Padrões](docs/patterns.md)** - Padrões de desenvolvimento
- **[Desenvolvimento](docs/development.md)** - Configuração e testes
- **[cmd/AGENTS.md](cmd/AGENTS.md)** - Guia de estrutura de comandos
- **[internal/AGENTS.md](internal/AGENTS.md)** - Estrutura de pacotes internos

## 📝 Licença

Este projeto é de código aberto e está disponível sob a Licença MIT.

## 🔗 Links

- [Repositório GitHub](https://github.com/Brennon-Oliveira/dev-cli)
- [Releases](https://github.com/Brennon-Oliveira/dev-cli/releases)
- [Especificação de Dev Container](https://containers.dev)

## ✨ Características

- ✅ Setup de container em um comando e integração com VS Code
- ✅ Suporte completo a WSL com conversão automática de caminhos
- ✅ Suporte para Docker e Podman
- ✅ Monitoramento de logs em tempo real
- ✅ Visibilidade de mapeamento de portas
- ✅ Autocompletar para bash, zsh e PowerShell
- ✅ Gerenciamento de configuração
- ✅ Funcionalidade experimental de auto-atualização
