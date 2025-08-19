package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"strconv"
	"strings"
)

type Sleep struct {
	Comando  string
	Mensagem *estruturas.Mensagem
}

func (instance Sleep) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.Comando)
	var err error

	if len(comandoSeparado) > 1 {
		instance.Mensagem.TempoEspera, err = strconv.Atoi(strings.TrimSpace(comandoSeparado[1]))
		if err != nil {
			resposta = err.Error()
		} else {
			resposta = "Sleep Atualizado - OK"
		}
	}
	return
}
