package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	. "d3c/server/helpers"
	"log"
)

type Select struct {
	ComandoCompleto   string
	AgenteSelecionado *string
	AgentesEmCampo    *[]estruturas.Mensagem
}

func (instance Select) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.ComandoCompleto)

	if len(comandoSeparado) > 1 {
		if AgenteCadastrado(&comandoSeparado[1], instance.AgentesEmCampo) {
			*instance.AgenteSelecionado = comandoSeparado[1]
		} else {
			log.Println("O Agente selecionado não está em campo.")
			log.Println("Para listar os AgentesEmCampo em campo use: show agentes.")
		}
	} else {
		*instance.AgenteSelecionado = ""
	}

	return ""
}
