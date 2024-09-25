package main

import (
	"EmailIndex/pkg/config"
	"EmailIndex/pkg/model"
	"EmailIndex/pkg/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

//var jSonFinal []string //array donde se guardará todos los correos en forma de objetos.
// List all folders

func main() {
	// Load configuration
	config.LoadConfig()

	// Start the CPU profiling
	start := time.Now()
	cpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()

	path := config.AppConfig.EmailDataPath
	fmt.Println("Indexando...")
	user_list := utils.ListAllFolders(path)

	// Create a channel to receive email data
	emailChan := make(chan model.Email, 1000)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start the indexing goroutine
	go model.BatchIndexData(emailChan)

	for _, user := range user_list {
		folders := utils.ListAllFolders(path + user)
		for _, folder := range folders {
			mail_files := utils.ListFiles(path + user + "/" + folder + "/")
			for _, mail_file := range mail_files {
				wg.Add(1)
				go func(user, folder, mail_file string) {
					defer wg.Done()
					fmt.Println("Indexing: " + user + "/" + folder + "/" + mail_file)
					sys_file, err := os.Open(path + user + "/" + folder + "/" + mail_file)
					if err != nil {
						log.Printf("Error opening file %s: %v", mail_file, err)
						return
					}
					defer sys_file.Close()

					email := model.ParseData(bufio.NewScanner(sys_file))
					emailChan <- email
				}(user, folder, mail_file)
			}
		}
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Close the email channel
	close(emailChan)
	//utils.JSONfinal(jSonFinal)
	fmt.Println("Indexing finished!!!!")
	// Calcular la duración del proceso
	duration := time.Since(start)
	fmt.Printf("El proceso tomó %v\n", duration)
	//Proceso de rendimiento de la aplicación
	runtime.GC()
	mem, err := os.Create("memory.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer mem.Close()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		log.Fatal(err)
	}
	////Fin proceso de rendimiento de la aplicación/////////////

}
