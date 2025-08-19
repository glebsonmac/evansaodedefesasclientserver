package helpers

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"os"
	"strings"
)

func RetornaNomeDoExecutavel() (nome string) {

	execPath := strings.Split(os.Args[0], "\\")
	nome = execPath[len(execPath)-1]

	return nome
}

func EvitaMultiplosAgentes() {
	nprocs := 0
	processos, _ := ps.Processes()
	executavel := RetornaNomeDoExecutavel()

	for proc := range processos {
		processo := processos[proc]
		procname := fmt.Sprintf("%s", processo.Executable())

		if procname == executavel {
			nprocs++
		}
	}

	if nprocs > 1 {
		os.Exit(1)
	}
}
