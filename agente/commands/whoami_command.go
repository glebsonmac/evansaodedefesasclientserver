package commands

import "os/user"

type Whoami struct {
}

func (instance Whoami) Executar() (resposta string) {
	usuario, err := user.Current()

	if err != nil {
		resposta = err.Error()
	} else {
		resposta = usuario.Username
	}
	return
}
