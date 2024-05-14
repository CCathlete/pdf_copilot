package pdfhandler

import (
	"fmt"
	"log"
	"maps"
	"os"
	"os/exec"
	"regexp"

	pdf2txt "github.com/heussd/pdftotext-go"
)

type ParasiteInfo map[string]string

func (p *ParasiteInfo) Init() {
	initialPattern := ParasiteInfo{
		"Scientific Name": "",
		"Common Name":     "",
		"Adult Size":      "",
		"Egg Size":        "",
		"Importance":      "",
		"Diagnosis":       "",
		"Treatment":       "",
		"Note":            "",
	}
	maps.Copy(*p, initialPattern)
}

func ConvertToText(pdfPath string) {
	cmd := exec.Command("bash", "-c", "pdftotext "+pdfPath)
	stdOutErr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(stdOutErr)
		log.Fatalf("Failed to run pdftotext using the path %s, the error was: %v\n", pdfPath, err)
	}
}

func ExtractParasitesInfo(txtPath string) (output []ParasiteInfo) {
	txtBytes, err := os.ReadFile(txtPath)
	if err != nil {
		log.Fatalf("Failed to read txt document: %v\n", err)
	}
	txtString := string(txtBytes)
	pattern := `((.*\n*)*)Fig` //Prototype pattern, change in the future.
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(txtString, -1)
	for _, match := range matches {
		pInfo := ParasiteInfo{}
		pInfo.Init()
		pInfo["Note"] = match
		output = append(output, pInfo)
	}
	return output
}

func ParsePdf(pdfPath string) []pdf2txt.PdfPage {
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		log.Fatalf("Failed to read PDF: %v\n", err)
	}

	/*
		pages = []PdfPage
		type PdfPage struct{
			Content string -> page text content
			Number int -> page number
		}
	*/
	pages, err := pdf2txt.Extract(pdfBytes)
	if err != nil {
		log.Fatalf("Failed to extract text from pages: %v\n", err)
	}

	return pages
}
