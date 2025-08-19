package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	. "d3c/server/helpers"
	"log"
	"os"
)

type Send struct {
	ComandoCompleto   string
	AgenteSelecionado *string
	AgentesEmCampo    *[]estruturas.Mensagem
}

func (instance Send) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.ComandoCompleto)

	if len(comandoSeparado) > 1 && *instance.AgenteSelecionado != "" {
		var erro error
		arquivoParaEnviar := &estruturas.Arquivo{}

		arquivoParaEnviar.Nome = comandoSeparado[1]
		arquivoParaEnviar.Conteudo, erro = os.ReadFile(arquivoParaEnviar.Nome)

		comandoSend := &estruturas.Commando{}
		comandoSend.Comando = comandoSeparado[0]
		comandoSend.Arquivo = *arquivoParaEnviar

		if erro != nil {
			log.Println("Erro ao abrir o arquivo: ", erro.Error())
		} else {
			posicaoDoAgenteNoArray := PosicaoAgenteEmCampo(*instance.AgenteSelecionado, *instance.AgentesEmCampo)

			(*instance.AgentesEmCampo)[posicaoDoAgenteNoArray].Comandos =
				append((*instance.AgentesEmCampo)[posicaoDoAgenteNoArray].Comandos, *comandoSend)
		}
	} else {
		log.Println("Especifique o arquivo a ser enviado.")
	}

	return ""
}
