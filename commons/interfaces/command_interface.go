package interfaces

type Command interface {
	Executar() (resposta string)
}
