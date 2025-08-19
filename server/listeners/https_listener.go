package listeners

import (
	"d3c/commons/estruturas"
	. "d3c/server/helpers"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	certFile = "server.crt"
	keyFile  = "server.key"
)

type httpsHandler struct{}

func (self *httpsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mensagem := &estruturas.Mensagem{}
	requestBytes, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Errorf("Erro ao receber requisicao: %v", err)
		return
	}

	err = json.Unmarshal(requestBytes, mensagem)
	if err != nil {
		fmt.Errorf("Erro ao tratar requisicao: %v", err)
		return
	}

	responseBody := TrataMensagem(mensagem, listaDeAgentes, agenteEmExecucao)
	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Errorf("Erro ao tratar requisicao: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(responseBytes)
	if err != nil {
		fmt.Errorf("Erro ao enviar resposta: %v", err)
		return
	}
}

func StartHttpsListener(port string, agentesEmCampo *[]estruturas.Mensagem, agenteSelecionado *string) {
	handler := &httpsHandler{}

	listaDeAgentes = agentesEmCampo
	agenteEmExecucao = agenteSelecionado

	err := http.ListenAndServeTLS(":"+port, certFile, keyFile, handler)

	if err != nil {
		fmt.Errorf("Erro ao iniciar o HTTPS listener: %v", err)
	}
}
