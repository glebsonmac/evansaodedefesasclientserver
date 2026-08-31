package connectors

import (
	"bytes"
	"d3c/commons/estruturas"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HttpConnector struct{}

func (self *HttpConnector) Execute(servidor, porta string, mensagem *estruturas.Mensagem) {
	client := &http.Client{Timeout: 15 * time.Second}

	bodyToSend, err := json.Marshal(*mensagem)
	if err != nil {
		return
	}

	request, err := http.NewRequest("POST", "http://"+servidor+":"+porta, bytes.NewBuffer(bodyToSend))
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-D3C-Agent", "1")

	serverResponse, err := client.Do(request)
	if err != nil {
		return
	}
	defer serverResponse.Body.Close()
	mensagem.Comandos = []estruturas.Commando{}

	bodyBytes, err := io.ReadAll(serverResponse.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bodyBytes, &mensagem); err != nil {
		return
	}
}
