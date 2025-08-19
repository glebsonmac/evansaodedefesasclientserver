package commands

import "os"

type Pwd struct {
}

func (instance Pwd) Executar() (resposta string) {
	resposta, err := os.Getwd()

	if err != nil {
		resposta = err.Error()
	}

	return
}
