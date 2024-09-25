package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ListAllFolders(folderName string) []string {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	listFolders := make([]string, 0, len(files)) // Pre-allocate slice
	for _, f := range files {
		if f.IsDir() {
			listFolders = append(listFolders, f.Name())
		}
	}
	return listFolders
}

// Lista cada uno de los archivos o correos
func ListFiles(folderName string) []string {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	fileNames := make([]string, 0, len(files)) // Pre-allocate slice
	for _, file := range files {
		if !file.IsDir() { // Only add files, not directories
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func JSONfinal(datos []string) {
	file, err := os.Create("jSonFinal.json")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Ensure file is closed

	// Use a buffered writer for better performance
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString(`{"Enron-email": [`)
	for i, dato := range datos {
		if i > 0 {
			writer.WriteString(",")
		}
		writer.WriteString(dato)
	}
	writer.WriteString("]}")

	fmt.Println("JSON File successfully created")
}
