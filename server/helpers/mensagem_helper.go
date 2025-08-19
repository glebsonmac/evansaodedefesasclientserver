package helpers

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"log"
)

func MensagemContemResposta(mensagem *estruturas.Mensagem) (contem bool) {
	contem = false

	for _, comando := range mensagem.Comandos {
		if len(comando.Resposta) > 0 {
			contem = true
		}
	}

	return contem
}

func TrataMensagem(mensagem *estruturas.Mensagem, agentesEmCampo *[]estruturas.Mensagem, agenteSelecionado *string) (mensagemParaOAgente estruturas.Mensagem) {
	if AgenteCadastrado(&mensagem.AgentID, agentesEmCampo) {
		if MensagemContemResposta(mensagem) {
			log.Println("Resposta do Host: ", mensagem.AgentHostname)
			// Exibir as respostas
			for indice, comando := range mensagem.Comandos {
				log.Println("Resposta do ComandoCompleto: ", comando.Comando)
				println(comando.Resposta)
				if helpers.SeparaComando(comando.Comando)[0] == "get" &&
					mensagem.Comandos[indice].Arquivo.Erro == false {
					SalvarArquivo(mensagem.Comandos[indice].Arquivo)
				}
			}
		}
		posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(*agenteSelecionado, *agentesEmCampo)

		// Enviar a lista de comandos enfileirados para o agente
		mensagemParaOAgente = (*agentesEmCampo)[posicaoDoAgenteNoArray]
		// Zera a lista de comandos do agente
		(*agentesEmCampo)[posicaoDoAgenteNoArray].Comandos = []estruturas.Commando{}
	} else {
		log.Println("Nova conexão! ID: ", mensagem.AgentID)
		(*agentesEmCampo) = append((*agentesEmCampo), *mensagem)
		mensagemParaOAgente = *mensagem
	}
	SetAlive(mensagem.AgentID, agentesEmCampo)
	return
}
