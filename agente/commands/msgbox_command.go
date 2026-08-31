package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Msgbox struct{}

func (instance Msgbox) Executar() (resposta string) {
	if runtime.GOOS != "windows" {
		return "Comando disponivel apenas no Windows."
	}

	tmpFile := filepath.Join(os.TempDir(), "aviso.txt")
	if err := os.WriteFile(tmpFile, []byte("Você foi hackeada.\n"), 0644); err != nil {
		return "Erro ao criar arquivo: " + err.Error()
	}

	if err := exec.Command("notepad.exe", tmpFile).Start(); err != nil {
		return "Erro ao abrir notepad: " + err.Error()
	}

	return "Notepad aberto com sucesso."
}
