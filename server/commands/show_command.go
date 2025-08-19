package commands

import (
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	. "d3c/server/helpers"
	"github.com/fatih/color"
	"log"
)

type Show struct {
	ComandoCompleto string
	AgentesEmCampo  *[]estruturas.Mensagem
}

func (instance Show) Executar() (resposta string) {
	comandoSeparado := helpers.SeparaComando(instance.ComandoCompleto)
	printVerde := color.New(color.FgGreen)
	printVermelho := color.New(color.FgRed)

	if len(comandoSeparado) > 1 {
		switch comandoSeparado[1] {
		case "agentes":
			for _, agente := range *instance.AgentesEmCampo {
				SetDead(agente.AgentID, instance.AgentesEmCampo)
				if len(comandoSeparado) > 2 && comandoSeparado[2] == "all" {
					if agente.AgentName != "" {
						if agente.AgenteDead {
							printVermelho.Println("Agente: " + agente.AgentName + " @ " + agente.AgentHostname)
						} else {
							printVerde.Println("Agente: " + agente.AgentName + " @ " + agente.AgentHostname)
						}
					} else {
						if agente.AgenteDead {
							printVermelho.Println("Agente: " + agente.AgentID + " @ " + agente.AgentHostname)
						} else {
							printVerde.Println("Agente: " + agente.AgentID + " @ " + agente.AgentHostname)
						}
					}
				} else {
					if !agente.AgenteDead {
						if agente.AgentName != "" {
							printVerde.Println("Agente: " + agente.AgentName + " @ " + agente.AgentHostname)
						} else {
							printVerde.Println("Agente: " + agente.AgentID + " @ " + agente.AgentHostname)
						}
					}
				}
			}
		default:
			log.Println("O parâmetro selecionado não existe.")
		}
	}

	return
}
