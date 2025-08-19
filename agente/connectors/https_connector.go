package connectors

import (
	"bytes"
	"crypto/tls"
	"d3c/commons/estruturas"
	"encoding/json"
	"io"
	"net/http"
)

type HttpsConnector struct {
}

func (self *HttpsConnector) Execute(servidor, porta string, mensagem *estruturas.Mensagem) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	bodyToSend, err := json.Marshal(*mensagem)
	if err != nil {
		return
	}

	request, err := http.NewRequest("POST", "https://"+servidor+":"+porta, bytes.NewBuffer(bodyToSend))
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/json")

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

	json.Unmarshal(bodyBytes, &mensagem)
}
