package helpers

import (
	"d3c/commons/estruturas"
	"log"
	"os"
)

func SalvarArquivo(arquivo estruturas.Arquivo) {
	err := os.WriteFile(arquivo.Nome, arquivo.Conteudo, 0644)
	if err != nil {
		log.Println("Erro ao salvar o arquivo recebido: ", err.Error())
	}
}
