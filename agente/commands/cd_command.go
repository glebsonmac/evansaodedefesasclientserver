package commands

import (
	"d3c/commons/helpers"
	"os"
	"strings"
)

type Cd struct {
	Comando string
}

func (instance Cd) Executar() (resposta string) {
	resposta = "Diretorio corrente alterado com sucesso!"
	comandoSeparado := helpers.SeparaComando(strings.TrimRight(instance.Comando, " "))

	if len(comandoSeparado) > 1 {
		err := os.Chdir(comandoSeparado[1])
		if err != nil {
			resposta = "Erro ao acessar o diretorio: " + comandoSeparado[1] + " -> " + err.Error()
		}
	}

	return
}
