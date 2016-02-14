package main

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const toolPath = "/Users/zat/Downloads/nand2tetris/tools/TextComparer.sh"

func TestMain(t *testing.T) {
	files, err := filepath.Glob("fixtures/*.jack")
	if err != nil {
		panic(err)
	}

	for _, jackFile := range files {
		testMainOutput(t, jackFile)
	}
}

func testMainOutput(t *testing.T, jackFile string) {
	xmlFile := regexp.MustCompile(`\.jack$`).ReplaceAllString(jackFile, ".xml")

	name := strings.Split(filepath.Base(jackFile), ".")[0]

	output, err := exec.Command("go", "run", "main.go", jackFile).CombinedOutput()
	if err != nil {
		t.Errorf("%s: %s %v", name, output, err)
		return
	}

	file, _ := ioutil.TempFile("", "")
	file.Write(output)

	output, err = exec.Command(toolPath, file.Name(), xmlFile).CombinedOutput()

	if err != nil {
		t.Errorf("%s: %s %v", name, output, err)
	}
}
