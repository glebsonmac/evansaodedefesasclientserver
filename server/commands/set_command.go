package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	. "d3c/server/helpers"
	"fmt"
)

type Set struct {
	ComandoCompleto   string
	AgenteSelecionado *string
	AgentesEmCampo    *[]estruturas.Mensagem
}

func (instance Set) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.ComandoCompleto)

	if len(comandoSeparado) > 2 && *instance.AgenteSelecionado != "" {
		switch comandoSeparado[1] {
		case "name":
			posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(*instance.AgenteSelecionado, *instance.AgentesEmCampo)
			(*instance.AgentesEmCampo)[posicaoDoAgenteNoArray].AgentName = comandoSeparado[2]
		default:
			fmt.Println("Opcao nao implementada para o comando `set`.")
		}
	} else if len(comandoSeparado) >= 2 && *instance.AgenteSelecionado == "" {
		fmt.Println("Selecione um agente para usar este comando.")
	} else if len(comandoSeparado) < 2 {
		fmt.Println("Defina a opcao para o comando.")
	}
	return ""
}
