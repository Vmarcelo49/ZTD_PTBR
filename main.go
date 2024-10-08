package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func getAITrad(line string) string {
	AILine, err := sendRequest(line)
	check(err)
	return AILine
}

func readFile(path string) []string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func isZTDString(str string) bool {
	var controle bool
	if strings.HasPrefix(str, "====") {
		controle = false
	}
	if strings.HasSuffix(str, "\\") {
		controle = true
	}
	return controle
}

type translationOBJ struct {
	//originalLine string
	aiTranslated string
	//translatedLine string
	index string
}

func newTrad(num string, ori string) translationOBJ {
	return translationOBJ{
		//originalLine: ori,
		index:        num,
		aiTranslated: getAITrad(ori),
	}

}

func trimmIndex(text string) string {
	newStr := strings.Replace(text, "====", "", -1)
	newStr = strings.Replace(newStr, " ", "", -1)
	return newStr
}

func buscarArquivosTXT() ([]string, error) {
	diretorio := "./text/"
	var arquivosTXT []string

	arquivos, err := os.ReadDir(diretorio)
	if err != nil {
		return nil, err
	}

	for _, arquivo := range arquivos {
		if strings.HasSuffix(arquivo.Name(), ".txt") {
			arquivosTXT = append(arquivosTXT, filepath.Join(diretorio, arquivo.Name()))
		}
	}
	fmt.Println("Encontrados os arquivos: ", arquivosTXT)
	return arquivosTXT, nil
}

func traduzir(arquivo string) {
	lines := readFile(arquivo)
	tradOBJS := []translationOBJ{}
	var lastIndex string
	var stringBuffer string
	var makingNewSTR bool
	for _, line := range lines {
		if makingNewSTR {
			if !isZTDString(line) {
				tradOBJ := newTrad(lastIndex, stringBuffer)
				tradOBJS = append(tradOBJS, tradOBJ)
				makingNewSTR = false
				stringBuffer = ""

			}
		}
		if !isZTDString(line) {
			lastIndex = trimmIndex(line)
		}
		if isZTDString(line) {
			trimmedStr := strings.Replace(line, "\\", "", -1)
			makingNewSTR = true
			stringBuffer = stringBuffer + " " + trimmedStr

		}
	}
	saveFile(tradOBJS, arquivo)

}

func main() {
	start := time.Now()
	arquivos, err := buscarArquivosTXT()
	check(err)
	for _, arquivo := range arquivos {
		traduzir(arquivo)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Tempo de execução: %s\n", elapsed)

}
