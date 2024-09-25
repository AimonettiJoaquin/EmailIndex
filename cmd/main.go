package main

import (
	"EmailIndex/pkg/config"
	"EmailIndex/pkg/model"
	"EmailIndex/pkg/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

func main() {
	// Load configuration settings
	config.LoadConfig()

	// Set up performance profiling
	start := time.Now()
	cpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()

	// Get the path where emails are stored
	path := config.AppConfig.EmailDataPath

	// List all user folders
	user_list := utils.ListAllFolders(path)

	// Create a channel to send emails through
	emailChan := make(chan model.Email, 1000)

	// Create a WaitGroup to keep track of running tasks
	var wg sync.WaitGroup

	// Start a goroutine to handle indexing emails
	go model.BatchIndexData(emailChan)

	// Loop through all users, folders, and email files
	for _, user := range user_list {
		folders := utils.ListAllFolders(path + user)
		for _, folder := range folders {
			mail_files := utils.ListFiles(path + user + "/" + folder + "/")
			for _, mail_file := range mail_files {
				// For each email file, start a new goroutine to process it
				wg.Add(1)
				go func(user, folder, mail_file string) {
					defer wg.Done()
					// Open the file, parse the email, and send it to the channel
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

	// Wait for all email processing to finish
	wg.Wait()
	// Close the email channel
	close(emailChan)

	// Print finish message and duration
	fmt.Println("Indexing finished!!!!")
	// Calculate the duration of the process
	duration := time.Since(start)
	fmt.Printf("The process took %v\n", duration)

	// Perform memory profiling
	mem, err := os.Create("memory.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer mem.Close()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		log.Fatal(err)
	}
	//go tool pprof -http=:8080 cpu.prof

}
