package helpers

func ValidaComando(comando string) (id int) {
	mapping := map[string]int{
		"cd":      1,
		"ls":      2,
		"ps":      3,
		"pwd":     4,
		"whoami":  5,
		"sleep":   6,
		"send":    7,
		"get":     8,
		"persist": 9,
	}

	id, _ = mapping[comando]

	return
}
