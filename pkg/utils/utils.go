package utils

import (
	"fmt"
	"log"
	"os"
)

func ListAllFolders(folder_name string) []string { //recibe como parámetro el folder "maildir".
	files, err := os.ReadDir(folder_name) //"ioutil.ReadDir" extrae todos los subfolders y los guarda en "files"
	if err != nil {
		log.Fatal(err)
	}
	var list_folders []string //array donde se guardarán las subcarpetas de "maildir"
	for _, f := range files {
		if f.IsDir() { //Si es un directorio
			list_folders = append(list_folders, f.Name()) //Guradmos el nombre de cada subfolder
		}

	}
	return list_folders
}

// Lista cada uno de los archivos o correos
func ListFiles(folder_name string) []string {
	files, err := os.ReadDir(folder_name)
	if err != nil {
		log.Fatal(err)
	}
	var files_names []string //array donde se guardarán los nombres de los archivos contenidos en las subcarpetas.
	for _, file := range files {
		files_names = append(files_names, file.Name())
	}
	return files_names
}

func JSONfinal(datos []string) {
	file, err := os.Create("jSonFinal.json")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	file.WriteString("{")
	file.WriteString(`"Enron-email"` + ": [")
	for index := range datos {
		file.WriteString(datos[index])
		if index == len(datos)-1 {
			file.WriteString("]")
			file.WriteString("}")
		} else {
			file.WriteString(",")
		}
	}
	file.Close()
	fmt.Println("JSON File successfully created")
}
