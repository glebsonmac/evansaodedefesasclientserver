package commands

import (
	"d3c/commons/helpers"
	"io/fs"
	"os"
	"strings"
)

type Ls struct {
	Comando string
}

func (instance Ls) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(strings.TrimRight(instance.Comando, " "))
	var arquivos []fs.DirEntry
	var err error

	if len(comandoSeparado) > 1 {
		arquivos, err = os.ReadDir(comandoSeparado[1])
	} else {
		arquivos, err = os.ReadDir(".")
	}

	if err != nil {
		resposta = err.Error()
	} else {
		for _, arquivo := range arquivos {
			resposta += arquivo.Name() + "\n"
		}
	}

	return
}
