package helpers

func ValidaComando(comando string) (id int) {
	mapping := map[string]int{
		"show":   1,
		"sleep":  2,
		"select": 3,
		"send":   4,
		"get":    5,
		"set":    6,
		"start":  7,
	}

	id, _ = mapping[comando]

	return
}
