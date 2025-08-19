package interfaces

import "d3c/commons/estruturas"

type Connector interface {
	Execute(servidor string, porta string, mensagem *estruturas.Mensagem)
}
