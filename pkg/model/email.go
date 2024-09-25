package model

import (
	"EmailIndex/pkg/config"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Estructura del json que contendrá cada correo.
type Email struct {
	ID                        int    `json:"ID"`
	Message_ID                string `json:"Message-ID"`
	Date                      string `json:"Date"`
	From                      string `json:"from"`
	To                        string `json:"to"`
	Subject                   string `json:"subject"`
	Mime_Version              string `json:"Mime-Version"`
	Content_Type              string `json:"Content-Type"`
	Content_Transfer_Encoding string `json:"Content-Transfer-Encoding"`
	X_From                    string `json:"X-From"`
	X_To                      string `json:"X-To"`
	X_cc                      string `json:"X-cc"`
	X_bcc                     string `json:"X-bcc"`
	X_Folder                  string `json:"X-Folder"`
	X_Origin                  string `json:"X-Origin"`
	X_FileName                string `json:"X-FileName"`
	Cc                        string `json:"Cc"`
	Body                      string `json:"Body"`
}

func ParseData(data_lines *bufio.Scanner) Email {
	var data Email
	var bodyStarted bool
	headerMap := map[string]*string{
		"Message-ID:":                &data.Message_ID,
		"Date:":                      &data.Date,
		"From:":                      &data.From,
		"To:":                        &data.To,
		"Subject:":                   &data.Subject,
		"Cc:":                        &data.Cc,
		"Mime-Version:":              &data.Mime_Version,
		"Content-Type:":              &data.Content_Type,
		"Content-Transfer-Encoding:": &data.Content_Transfer_Encoding,
		"X-From:":                    &data.X_From,
		"X-To:":                      &data.X_To,
		"X-cc:":                      &data.X_cc,
		"X-bcc:":                     &data.X_bcc,
		"X-Folder:":                  &data.X_Folder,
		"X-Origin:":                  &data.X_Origin,
		"X-FileName:":                &data.X_FileName,
	}

	for data_lines.Scan() {
		line := data_lines.Text()
		if !bodyStarted {
			for prefix, field := range headerMap {
				if strings.HasPrefix(line, prefix) {
					*field = strings.TrimSpace(line[len(prefix):])
					break
				}
			}
			if line == "" {
				bodyStarted = true
			}
		} else {
			data.Body += line + "\n"
		}
	}
	return data
}

// BatchIndexData handles batch indexing of emails
func BatchIndexData(emailChan <-chan Email) {
	batchSize := config.AppConfig.BatchSize
	var batch []Email
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case email, ok := <-emailChan:
			if !ok {
				// Channel closed, index remaining emails
				if len(batch) > 0 {
					indexBatch(batch)
				}
				return
			}
			batch = append(batch, email)
			if len(batch) >= batchSize {
				indexBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				indexBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// Esta función realiza una petición y envía los datos.
func indexBatch(batch []Email) {
	user := config.AppConfig.ZincUser
	password := config.AppConfig.ZincPassword
	auth := user + ":" + password
	bas64encoded_creds := base64.StdEncoding.EncodeToString([]byte(auth))
	index := config.AppConfig.ZincIndex
	zinc_host := config.AppConfig.ZincHost
	zinc_url := zinc_host + "/api/" + index + "/_bulk"

	var bulkRequestBody bytes.Buffer
	for _, email := range batch {
		// Add metadata line
		metadata := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
			},
		}
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			log.Printf("Error marshaling metadata: %v", err)
			continue
		}
		bulkRequestBody.Write(metadataJSON)
		bulkRequestBody.WriteString("\n")

		// Add document line
		documentJSON, err := json.Marshal(email)
		if err != nil {
			log.Printf("Error marshaling email: %v", err)
			continue
		}
		bulkRequestBody.Write(documentJSON)
		bulkRequestBody.WriteString("\n")
	}

	req, err := http.NewRequest("POST", zinc_url, &bulkRequestBody)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+bas64encoded_creds)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error indexing batch. Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))
	} else {
		log.Printf("Successfully indexed batch of %d documents", len(batch))
	}
}
