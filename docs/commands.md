# Criando Novos Comandos

O Dev CLI utiliza a biblioteca [Cobra](https://github.com/spf13/cobra) para gerenciar seus comandos.

## Estrutura de um Comando

Para criar um novo comando, adicione um arquivo `.go` no diretório `cmd/`. Por exemplo, `cmd/exemplo.go`:

```go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var exemploCmd = &cobra.Command{
	Use:   "exemplo [argumento]",
	Short: "Breve descrição do comando",
	Long:  "Descrição detalhada sobre o que o comando faz.",
	Args:  cobra.ExactArgs(1), // Define a obrigatoriedade de argumentos
	RunE: func(cmd *cobra.Command, args []string) error {
		// Lógica do comando
		fmt.Printf("Executando exemplo com: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exemploCmd)
}
```

## Padrões para Comandos

1. **Validação de Argumentos**: Use `Args: cobra.MinimumNArgs(x)` ou funções similares para garantir que o usuário passe os parâmetros necessários.
2. **Uso de RunE**: Prefira `RunE` em vez de `Run`. Isso permite retornar erros diretamente para o Cobra, que os exibirá de forma padronizada.
3. **Delegar Lógica**: Evite colocar lógica de negócio pesada dentro do pacote `cmd`. Em vez disso, chame funções em `internal/container` ou outros pacotes internos.
4. **Tratamento de Caminhos**: Sempre use `container.GetAbsPath(path)` para normalizar os caminhos passados pelo usuário.
