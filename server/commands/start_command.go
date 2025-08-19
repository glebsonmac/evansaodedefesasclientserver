package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	. "d3c/server/listeners"
)

type Start struct {
	ComandoCompleto   string
	AgenteSelecionado *string
	AgentesEmCampo    *[]estruturas.Mensagem
}

func (instance Start) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.ComandoCompleto)

	if len(comandoSeparado) > 2 {
		switch comandoSeparado[1] {
		case "raw":
			go StartRawListener(comandoSeparado[2], instance.AgentesEmCampo, instance.AgenteSelecionado)
		case "https":
			go StartHttpsListener(comandoSeparado[2], instance.AgentesEmCampo, instance.AgenteSelecionado)
		}
	}

	return ""
}
