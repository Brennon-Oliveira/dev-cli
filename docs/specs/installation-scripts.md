# Spec: Instalação e Atualização Unificada (dev-cli)

## Resumo
Scripts para simplificar o onboarding do `dev-cli`, agregando funções de "Instalação do zero" e "Atualização" em uma execução inteligente, para Windows (`.ps1`) e Unix-like (`.sh`). Além de centralizar a instalação, padronizaremos a nomenclatura do binário para `dev-cli` em todo o ecossistema (removendo referências soltas a apenas `dev`).

## Objetivos e Restrições
1. **Padronização do Nome:** Garantir que em todos os builds (CI/CD) e referências de auto-completion, o binário se chame `dev-cli`. O `cmd/root.go` deverá usar `Use: "dev-cli"`.
2. **Fim do Snap:** A instalação via snap apresentou problemas. O snap deve ser removido de toda a documentação, CI/CD e configuração.
3. **Instalação Automática e Manual:** O `README.md` e `README.pt.md` apresentarão a instalação "One-Liner" (Automática) como padrão, mas manterão uma seção de Instalação Manual como alternativa.
4. **Comportamento do `dev-cli update`:** O comando `update` interno do Go **NÃO** deve alterar o próprio binário em execução. Em vez disso, se houver uma nova versão, ele deve informar o usuário e sugerir o comando exato (`curl` ou `irm`) para atualizar a ferramenta.

## Fluxos de Execução dos Scripts (Cenários)

1. **Nova Instalação (Unix/Windows):**
   - Script é baixado e executado.
   - Detecta ausência do `dev-cli` no PATH.
   - Consulta API do GitHub por `releases/latest`.
   - Baixa o artefato, extrai, move para o destino (`/usr/local/bin` no Linux/macOS ou `$env:USERPROFILE\.dev-cli\bin` no Windows).
   - Ajusta o PATH do usuário, se necessário.
   - Mostra mensagem de sucesso.

2. **Atualização (Unix/Windows):**
   - Detecta `dev-cli` no PATH.
   - Pega versão local via `dev-cli --version`.
   - Consulta API do GitHub pela versão mais recente.
   - Se `local == remota`: Exibe "Already up to date." e finaliza.
   - Se `local != remota`: Pergunta `[Y/n]` para atualizar (se o usuário quiser). Se `Y`, faz o overwrite do executável atual.

3. **Fallback e Tratamento de Erros:**
   - Exibir erros amigáveis se houver limite de taxa na API do GitHub.
   - Exibir alerta se a arquitetura do sistema não for suportada.
   - Se houver falha de permissão ao mover o binário (ex. `/usr/local/bin`), sugerir que o usuário rode o comando final com permissão de superusuário (sudo) ou tentar rodar o próprio script em modo sudo.

## Todo-List de Implementação

- [ ] **Ajustes no Core e CI/CD:**
  - Em `cmd/root.go`, alterar `Use: "dev"` para `Use: "dev-cli"` e revisar os exemplos contidos na CLI.
  - No `.github/workflows/build.yml`, remover o job `snap-amd64` e remover a dependência desse job no release.
  - Excluir o diretório `snap/` do repositório.

- [ ] **Criar `/install.sh` (Bash/Zsh):**
  - Obter a tag mais recente através de chamadas à API do GitHub (usando `curl` ou `wget`).
  - Identificar OS (`uname -s`) e Arquitetura (`uname -m`).
  - Lógica para download de `dev-[OS]-[ARCH].tar.gz` e extração.
  - Validação de versão atual instalada para questionar sobre atualização (`[Y/n]`).

- [ ] **Criar `/install.ps1` (PowerShell):**
  - Obter tag mais recente usando `Invoke-RestMethod`.
  - Usar `$env:PROCESSOR_ARCHITECTURE` para definir a arquitetura no Windows.
  - Baixar `dev-windows-amd64.zip`, usar `Expand-Archive`, e copiar `dev-cli.exe` para o destino correto (`$env:USERPROFILE\.dev-cli\bin`).
  - Adicionar o diretório de destino ao `PATH` do `User`.

- [ ] **Atualizar o `cmd/update.go`:**
  - Alterar a mensagem de saída. Quando houver nova versão, imprimir: 
    * "Nova versão encontrada: [VERSAO_NOVA]".
    * "Para atualizar, execute:"
    * `curl -fsSL https://raw.githubusercontent.com/Brennon-Oliveira/dev-cli/main/install.sh | bash` (Se for Unix)
    * `irm https://raw.githubusercontent.com/Brennon-Oliveira/dev-cli/main/install.ps1 | iex` (Se for Windows)

- [ ] **Atualizar Documentações (`README.md` / `README.pt.md`):**
  - Remover qualquer menção ao `snap`.
  - Criar seção "Instalação Rápida" usando o comando de script.
  - Manter/Reescrever a seção "Instalação Manual" com as instruções de download e extração via linha de comando padrão.

## Exemplos de Código Relevantes

### Buscar versão e confirmar update no Bash (`install.sh`)
```bash
LATEST_TAG=$(curl -sL https://api.github.com/repos/Brennon-Oliveira/dev-cli/releases/latest | grep '"tag_name":' | cut -d'"' -f4)

if command -v dev-cli &> /dev/null; then
    CURRENT_VERSION=$(dev-cli --version | awk '{print $3}')
    if [ "$CURRENT_VERSION" = "$LATEST_TAG" ]; then
        echo "✅ dev-cli is already up to date ($CURRENT_VERSION)"
        exit 0
    else
        read -p "Update dev-cli from $CURRENT_VERSION to $LATEST_TAG? [Y/n] " confirm
        if [[ $confirm == [nN] || $confirm == [nN][oO] ]]; then
            echo "Update cancelled."
            exit 0
        fi
    fi
fi
```

### Prompt de Update no `cmd/update.go`
```go
if runtime.GOOS == "windows" {
    logger.Info("Para atualizar, execute no PowerShell:\n  irm https://raw.githubusercontent.com/Brennon-Oliveira/dev-cli/main/install.ps1 | iex")
} else {
    logger.Info("Para atualizar, execute no terminal:\n  curl -fsSL https://raw.githubusercontent.com/Brennon-Oliveira/dev-cli/main/install.sh | bash")
}
```
