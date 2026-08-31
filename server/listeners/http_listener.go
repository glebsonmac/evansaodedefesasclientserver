package listeners

import (
	"d3c/commons/estruturas"
	. "d3c/server/helpers"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type httpCombinedHandler struct {
	agentesEmCampo    *[]estruturas.Mensagem
	agenteSelecionado *string
	webMux            http.Handler
}

func (h *httpCombinedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-D3C-Agent") == "1" {
		h.handleAgent(w, r)
		return
	}
	h.webMux.ServeHTTP(w, r)
}

func (h *httpCombinedHandler) handleAgent(w http.ResponseWriter, r *http.Request) {
	mensagem := &estruturas.Mensagem{}
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Erro ao receber agente: %v", err)
		return
	}

	if err = json.Unmarshal(requestBytes, mensagem); err != nil {
		log.Printf("Erro ao deserializar agente: %v", err)
		return
	}

	responseBody := TrataMensagem(mensagem, h.agentesEmCampo, h.agenteSelecionado)
	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

func StartHttpListener(port string, agentesEmCampo *[]estruturas.Mensagem, agenteSelecionado *string) {
	handler := &httpCombinedHandler{
		agentesEmCampo:    agentesEmCampo,
		agenteSelecionado: agenteSelecionado,
		webMux:            newWebMux(agentesEmCampo, agenteSelecionado),
	}
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Printf("Erro no HTTP listener: %v", err)
	}
}
