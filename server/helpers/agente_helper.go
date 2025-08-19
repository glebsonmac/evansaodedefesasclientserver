package helpers

import (
	"d3c/commons/estruturas"
	"time"
)

func PosicaoAgenteEmCampo(agenteId string, agentesEmCampo []estruturas.Mensagem) (posicao int) {
	for indice, agente := range agentesEmCampo {
		if agenteId == agente.AgentID || agenteId == agente.AgentName {
			posicao = indice
		}
	}
	return posicao
}

func AgenteCadastrado(agentID *string, agentesEmCampo *[]estruturas.Mensagem) (cadastrado bool) {
	cadastrado = false

	for _, agente := range *agentesEmCampo {
		if agente.AgentID == *agentID || agente.AgentName == *agentID {
			cadastrado = true
		}
	}

	return cadastrado
}

func SetAlive(agentID string, agentesEmCampo *[]estruturas.Mensagem) {
	posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(agentID, *agentesEmCampo)
	(*agentesEmCampo)[posicaoDoAgenteNoArray].LastSeen = time.Now()
	(*agentesEmCampo)[posicaoDoAgenteNoArray].AgenteDead = false
}

func isDead(agentID string, agentesEmCampo *[]estruturas.Mensagem) (dead bool) {
	posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(agentID, *agentesEmCampo)
	return (*agentesEmCampo)[posicaoDoAgenteNoArray].AgenteDead
}

func SetDead(agentID string, agentesEmCampo *[]estruturas.Mensagem) {
	posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(agentID, *agentesEmCampo)
	tempoDeSleepDoAgente := (*agentesEmCampo)[posicaoDoAgenteNoArray].TempoEspera
	maxWaitTime := (*agentesEmCampo)[posicaoDoAgenteNoArray].LastSeen.Add(time.Duration(3*tempoDeSleepDoAgente) * time.Second)

	if time.Now().After(maxWaitTime) {
		(*agentesEmCampo)[posicaoDoAgenteNoArray].AgenteDead = true
	}

}
