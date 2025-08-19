package estruturas

import "time"

type Mensagem struct {
	AgentID       string
	AgentName     string
	LastSeen      time.Time
	AgenteDead    bool
	AgentHostname string
	AgentCWD      string
	TempoEspera   int
	Comandos      []Commando
}
