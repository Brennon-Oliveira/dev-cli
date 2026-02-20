# Dev CLI

Uma interface de linha de comando para gerenciamento de ciclo de vida de Dev Containers e integração nativa com o VS Code em modo *detached*.

## 📥 Instalação

Acesse os binários pré-compilados na [página de Releases](https://github.com/Brennon-Oliveira/dev-cli/releases/latest).

### Linux e macOS (Compatível nativamente com WSL)

Substitua `[OS]` por `linux` ou `macos` e `[ARCH]` por `amd64` ou `arm64` de acordo com a sua arquitetura.

```bash
# 1. Baixe o artefato compactado
curl -LO [https://github.com/Brennon-Oliveira/dev-cli/releases/latest/download/dev-](https://github.com/Brennon-Oliveira/dev-cli/releases/latest/download/dev-)[OS]-[ARCH].tar.gz

# 2. Descompacte o arquivo
tar -xzf dev-[OS]-[ARCH].tar.gz

# 3. Mova o executável para o PATH do sistema
sudo mv dev /usr/local/bin/

# 4. Limpe o arquivo baixado
rm dev-[OS]-[ARCH].tar.gz

# 5. Valide a instalação
dev --help

```

### Windows

1. Baixe o arquivo `dev-windows-amd64.zip` na [última release](https://www.google.com/url?sa=E&source=gmail&q=https://github.com/Brennon-Oliveira/dev-cli/releases/latest).
2. Extraia o conteúdo do `.zip`.
3. Mova o binário `dev.exe` para um diretório seguro (ex: `C:\Ferramentas\bin`).
4. Adicione este diretório à variável de ambiente `PATH` do Windows.
5. Valide a instalação executando `dev --help` no seu terminal preferido.

## 🛠️ Comandos

### Inicialização e Ciclo de Vida

* `dev run [caminho]`: (Recomendado) Provisiona o container e imediatamente abre o VS Code no diretório mapeado.
* `dev up [caminho]`: Provisiona e inicia o dev container em segundo plano, sem abrir o editor.
* `dev open [caminho]`: Abre o VS Code diretamente conectado ao dev container já em execução, resolvendo dinamicamente o `workspaceFolder` do `devcontainer.json`.
* `dev kill [caminho]`: Localiza e encerra instantaneamente o processo do container atrelado ao workspace alvo.

### Interação com o Ambiente

* `dev shell`: Injeta um shell interativo (`zsh`, `bash` ou `sh`) diretamente dentro do container ativo do diretório atual.
* `dev exec [comando]`: Repassa comandos e parâmetros arbitrários para serem executados no contexto isolado do container (ex: `dev exec npm run build`).

### Monitoramento e Diagnóstico

* `dev list`: Retorna a lista de todos os dev containers em execução no host local.
* `dev logs [-f]`: Exibe a saída padrão do container. Use a flag `-f` para acompanhamento em tempo real (*tail*).
* `dev ports`: Lista todos os mapeamentos de rede e portas expostas ativas entre o host e o container atual.

### Manutenção

* `dev clean`: Realiza a liberação de recursos do Docker, removendo containers parados e redes órfãs geradas pela extensão.

## ⚙️ Casos de Uso

* **Onboarding Imediato:** Após clonar um projeto, não é necessário abrir o VS Code, localizar a pasta e clicar em "Reopen in Container". Basta rodar `dev run` na raiz do repositório pelo terminal. A CLI resolve o build do Docker, lida com a interoperabilidade de caminhos (caso esteja usando WSL) e injeta o `code` na estrutura final.
* **Execução Headless:** Se você precisa apenas rodar testes ou compilar artefatos em um ambiente padronizado, utilize `dev up` para subir a infraestrutura invisível e `dev exec` para acionar as rotinas, consumindo menos memória do sistema host por não instanciar o Electron.
* **Resolução Avançada de Caminhos:** O projeto analisa as configurações do `devcontainer.json` nativamente através de regex, garantindo que o editor acesse a pasta raiz real (`workspaceFolder`), lidando automaticamente com fallbacks, caminhos curtos e bugs de parse de URI do VS Code.
