package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const mainName = "../index.html"

func HandleMain(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	file, err := os.ReadFile(mainName)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("No Content"))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("error in writing answer to client: %v", err)
		return
	}
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	convertedData, err := service.Translate(string(fileBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("Conversion error: %v", err), http.StatusInternalServerError)
		return
	}

	outputFilename := fmt.Sprintf("result_%s%s",
		time.Now().UTC().Format("2006_01_02_15_04_05"),
		filepath.Ext(header.Filename))

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		http.Error(w, "Error creating result file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	if _, err := outputFile.WriteString(convertedData); err != nil {
		http.Error(w, "Error writing to result file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%s", convertedData)
}
