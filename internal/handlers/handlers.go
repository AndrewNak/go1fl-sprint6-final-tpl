package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go1fl-sprint6-final-tpl/internal/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	paths := []string{"index.html", "../index.html", "../index.html"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			http.ServeFile(w, r, path)
			return
		}
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	log.Printf("Uploaded file: %s, size: %d bytes", handler.Filename, handler.Size)

	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file content: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	converted, err := service.AutoDetectAndConvert(string(content))
	if err != nil {
		log.Printf("Error converting content: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Сохраняем файл в корневую директорию проекта
	timestamp := time.Now().UTC().Format("20060102_150405")
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".txt"
	}

	var prefix string
	if service.IsMorseCode(string(content)) {
		prefix = "morse_to_text_"
	} else {
		prefix = "text_to_morse_"
	}

	filename := prefix + timestamp + ext

	original := string(content)
	
	// Сохраняем результат в новый файл
	if err := os.WriteFile(filename, []byte(converted), 0644); err != nil {
		log.Printf("Error writing output file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("Created output file: %s", filename)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	response := "Исходный текст:\n" + original + "\n\n" +
		"Результат конвертации:\n" + converted + "\n\n" +
		"Файл сохранен как: " + filename

	w.Write([]byte(response))
}
