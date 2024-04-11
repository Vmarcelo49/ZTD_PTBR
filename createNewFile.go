package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func saveFile(linhas []translationOBJ, filename string) {
	nomeArquivo := filepath.Base(filename)
	arquivo, err := os.Create("translatedText/" + nomeArquivo)
	if err != nil {
		//odeio como essa bosta reporta erros nativamente
		fmt.Println("Erro na criação do arquivo:", nomeArquivo, arquivo)
		panic(err)
	}
	defer arquivo.Close()
	for _, linha := range linhas {
		arquivo.WriteString("==== " + linha.index + " ====\n")
		linhaQuebrada := quebraLinha(linha.aiTranslated)
		arquivo.WriteString(linhaQuebrada + "\\\\\n")
	}
	fmt.Println("Terminado arquivo ", nomeArquivo)

}

func quebraLinha(str string) string {
	if len(str) <= 50 {
		return str
	}

	palavras := strings.Fields(str)
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
