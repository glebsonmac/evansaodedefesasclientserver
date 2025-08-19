package helpers

import "strings"

func SeparaComando(comandoCompleto string) (comandoSeparado []string) {
	comandoSeparado = strings.Split(comandoCompleto, " ")

	return comandoSeparado
}
