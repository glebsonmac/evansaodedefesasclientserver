package listeners

import (
	"d3c/commons/estruturas"
	. "d3c/server/helpers"
	"encoding/gob"
	"log"
	"net"
)

func StartRawListener(port string, agentesEmCampo *[]estruturas.Mensagem, agenteSelecionado *string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatal("Erro ao iniciar o listener: ", err.Error())
	} else {
		for {
			canal, err := listener.Accept()
			defer canal.Close()

			if err != nil {
				log.Println("Erro em um novo canal: ", err.Error())
				return
			}
			mensagem := &estruturas.Mensagem{}
			gob.NewDecoder(canal).Decode(mensagem)

			listaDeAgentes = agentesEmCampo
			agenteEmExecucao = agenteSelecionado

			responseBody := TrataMensagem(mensagem, listaDeAgentes, agenteEmExecucao)

			gob.NewEncoder(canal).Encode(responseBody)
		}
	}
}
