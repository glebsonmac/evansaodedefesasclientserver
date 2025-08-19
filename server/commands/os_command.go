package commands

import (
	"d3c/commons/estruturas"
	. "d3c/server/helpers"
	"log"
)

type Default struct {
	ComandoCompleto   string
	AgenteSelecionado *string
	AgentesEmCampo    *[]estruturas.Mensagem
}

func (instance Default) Executar() (resposta string) {

	if *instance.AgenteSelecionado != "" {
		comando := &estruturas.Commando{}
		comando.Comando = instance.ComandoCompleto

		posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(*instance.AgenteSelecionado, *instance.AgentesEmCampo)

		(*instance.AgentesEmCampo)[posicaoDoAgenteNoArray].Comandos =
			append((*instance.AgentesEmCampo)[posicaoDoAgenteNoArray].Comandos, *comando)
	} else {
		log.Println("O comando digitado não existe!")
	}
	return ""
}
