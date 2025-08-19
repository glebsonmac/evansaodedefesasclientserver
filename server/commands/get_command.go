package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	. "d3c/server/helpers"
	"log"
)

type Get struct {
	ComandoCompleto   string
	AgenteSelecionado *string
	AgentesEmCampo    *[]estruturas.Mensagem
}

func (instance Get) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.ComandoCompleto)

	if len(comandoSeparado) > 1 && *instance.AgenteSelecionado != "" {
		comandoSend := &estruturas.Commando{}
		comandoSend.Comando = instance.ComandoCompleto

		posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(*instance.AgenteSelecionado, *instance.AgentesEmCampo)

		(*instance.AgentesEmCampo)[posicaoDoAgenteNoArray].Comandos =
			append((*instance.AgentesEmCampo)[posicaoDoAgenteNoArray].Comandos, *comandoSend)

	} else {
		log.Println("Especifique o arquivo que deseja copiar.")
	}

	return ""
}
