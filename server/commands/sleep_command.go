package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	. "d3c/server/helpers"
	"log"
	"strconv"
)

type Sleep struct {
	ComandoCompleto string
	AgenteAlvo      string
	AgentesEmCampo  *[]estruturas.Mensagem
}

func (instance Sleep) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.ComandoCompleto)

	// TODO: separar a mensagem de erro
	if len(comandoSeparado) > 1 && instance.AgenteAlvo != "" {
		comandoSend := &estruturas.Commando{}
		comandoSend.Comando = instance.ComandoCompleto

		indiceDoAgenteNaLista := PosicaoAgenteEmCampo(instance.AgenteAlvo, *instance.AgentesEmCampo)
		(*instance.AgentesEmCampo)[indiceDoAgenteNaLista].TempoEspera, _ = strconv.Atoi(comandoSeparado[1])

		comandoParaOAgenteEmCampo := (*instance.AgentesEmCampo)[indiceDoAgenteNaLista].Comandos
		(*instance.AgentesEmCampo)[indiceDoAgenteNaLista].Comandos = append(comandoParaOAgenteEmCampo, *comandoSend)
	} else {
		log.Println("Escolha por quantos segundos o Agente deve dormir.")
	}

	return
}
