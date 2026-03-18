# Constants Package

## Responsabilidade

Centralizar todas as constantes do projeto para evitar strings mágicas e facilitar manutenção.

## Padrões

- Nenhuma string mágica no código - sempre usar constantes deste pacote
- Constantes em UPPER_SNAKE_CASE
- Agrupar constantes relacionadas logicamente
- Preferir `const` sobre `var` quando o valor é imutável

## Adicionando Nova Constante

1. Identificar o grupo lógico (Tool, Label, Config, etc.)
2. Adicionar ao bloco `const` apropriado
3. Se for uma lista de valores válidos, criar um slice `Valid*`
4. Atualizar referências no código para usar a constante

## Uso

```go
import "github.com/Brennon-Oliveira/dev-cli/internal/constants"

tool := constants.ToolDocker
label := constants.LabelDevContainerFolder
```
