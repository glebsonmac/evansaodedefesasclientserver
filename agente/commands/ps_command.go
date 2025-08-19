package commands

import (
	"fmt"
	"github.com/mitchellh/go-ps"
)

type Ps struct {
}

func (instance Ps) Executar() (resposta string) {
	listaDeProcessos, err := ps.Processes()

	if err != nil {
		resposta = err.Error()
	} else {
		resposta = fmt.Sprintf("1 - Parant ID, 2 - Process ID, 3 - Executable Name \n")
		for _, processo := range listaDeProcessos {
			resposta += fmt.Sprintf("\t%d -> %d -> %s \n", processo.PPid(), processo.Pid(), processo.Executable())
		}
	}

	return
}
