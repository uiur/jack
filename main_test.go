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
	jackFile := "fixtures/Main.jack"
	testMainOutput(t, jackFile)
}

func testMainOutput(t *testing.T, jackFile string) {
	xmlFile := regexp.MustCompile(`\.jack$`).ReplaceAllString(jackFile, ".xml")

	name := strings.Split(filepath.Base(jackFile), ".")[0]

	output, err := exec.Command("go", "run", "main.go", jackFile).Output()
	if err != nil {
		t.Error(err)
	}

	file, _ := ioutil.TempFile("", "")
	file.Write(output)

	output, err = exec.Command(toolPath, file.Name(), xmlFile).CombinedOutput()

	if err != nil {
		t.Errorf("%s: %v %v", name, output, err)
	}
}
