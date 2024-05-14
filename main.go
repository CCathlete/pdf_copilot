package main

import (
	"errors"
	"fmt"
	"os"
	dbH "pdf-to-postgres/dbHandler"
	pdf2txt "pdf-to-postgres/pdfHandler"
	ymlH "pdf-to-postgres/yamlHandler"
	"strings"
)

func main() {
	docPath := "Parasitology/Parasitology_book.pdf"
	txtPath := strings.Replace(docPath, "pdf", "txt", -1) // -1 means all instances.
	if _, err := os.Stat(txtPath); errors.Is(err, os.ErrNotExist) {
		pdf2txt.ConvertToText(docPath)
		fmt.Println("PDF doc was converted to text at path: " + txtPath)
		// fmt.Printf("The info: %v\n", parasites)
	} else {
		fmt.Println("The text file exists.")
	}
	parasites := pdf2txt.ExtractParasitesInfo(txtPath)

	yamlPath := "config.yaml"
	configMap := ymlH.ParseYaml(yamlPath)
	fmt.Printf("Our config is: %v\n", configMap)

	dbName := ymlH.GetDbName(yamlPath)
	fmt.Printf("Our DB name is: %s\n", dbName)

	dbInfo := configMap["Database"].(map[interface{}]interface{})
	dbPointer := dbH.DbInit(dbInfo)
	for _, parasiteInfo := range parasites {
		dbH.AddToTable(dbPointer, dbName, parasiteInfo)
	}
}

/*
In case we want to use the golang wrapper of poppler (pdftotext). It seemed to work better from the cli.
	pages := pdf2txt.ParsePdf(docPath)
	file, err := os.OpenFile(txtPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close the txt file: %v.", err)
		}
	}()
	if err != nil {
		log.Fatalf("Failed to append or create the txt file: %v.", err)
	}

	for _, page := range pages {
		_, err := file.WriteString(page.Content + "\n")
		if err != nil {
			log.Fatalf("Failed to write to the txt file: %v.", err)
		}
	}
*/
