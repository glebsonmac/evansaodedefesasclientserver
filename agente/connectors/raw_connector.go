package connectors

import (
	"d3c/commons/estruturas"
	"encoding/gob"
	"net"
)

type RawConnector struct {
}

func (self *RawConnector) Execute(servidor, porta string, mensagem *estruturas.Mensagem) {
	canal, _ := net.Dial("tcp", servidor+":"+porta)
	defer canal.Close()

	gob.NewEncoder(canal).Encode(*mensagem)
	mensagem.Comandos = []estruturas.Commando{}

	gob.NewDecoder(canal).Decode(&mensagem)
}
