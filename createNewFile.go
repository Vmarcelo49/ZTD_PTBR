package main

import (
	"os"
	"log"
	"path/filepath"
	"strings"
)

func saveFile(lines []translationOBJ, filename string) {
	nomeArquivo := filepath.Base(filename)
	arquivo, err := os.Create("translatedText/" + nomeArquivo)
	if err != nil {
		log.Panic("Erro na criação do arquivo:", nomeArquivo)
	}
	defer arquivo.Close()
	for _, linha := range lines {
		arquivo.WriteString("==== " + linha.index + " ====\n")
		linhaQuebrada := quebraLinha(linha.aiTranslated)
		arquivo.WriteString(linhaQuebrada + "\\\\\n")
	}
	log.Println("Terminado arquivo ", nomeArquivo)

}

func quebraLinha(linha string) string {
	if len(linha) <= 50 {
		return str
	}

	palavras := strings.Fields(linha)
	var resultado string
	//Não lembro de onde peguei essa parte, mas caraca como funciona bem
	contador := 0
	for _, palavra := range palavras {
		if contador+len(palavra) > 50 {
			resultado += "\\\\\n" + palavra
			contador = len(palavra) + 1
		} else {
			if contador > 0 {
				resultado += " "
			}
			resultado += palavra
			contador += len(palavra) + 1
		}
	}

	return resultado
}
