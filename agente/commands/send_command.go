package commands

import (
	"d3c/commons/estruturas"
	"os"
)

type Send struct {
	Arquivo estruturas.Arquivo
}

func (instance Send) Executar() (resposta string) {
	resposta = "Arquivo enviado com sucesso!"

	err := os.WriteFile(instance.Arquivo.Nome, instance.Arquivo.Conteudo, 0644)

	if err != nil {
		resposta = "Erro ao salvar arquivo no destino: " + err.Error()
	}

	return resposta
}
