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

	// Пробуем разные пути для index.html
	paths := []string{
		//"index.html",    // если запуск из корня
		//"../index.html", // если запуск из cmd
		"./index.html", // текущая директория
	}

	var html []byte
	var err error

	for _, path := range paths {
		html, err = os.ReadFile(path)
		if err == nil {
			log.Printf("Found index.html at: %s", path)
			break
		}
	}

	if err != nil {
		log.Printf("Error reading index.html: %v", err)
		http.Error(w, "Internal server error - index.html not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
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

	// Пробуем сохранить в разные места
	filename := prefix + timestamp + ext
	savePaths := []string{
		filename,         // текущая директория (cmd)
		"../" + filename, // корневая директория проекта
		"/home/andrew/practice/go1fl-sprint6-final-tpl/go1fl-sprint6-final-tpl/" + filename, // абсолютный путь
	}

	var outputFile *os.File
	var savedPath string

	for _, path := range savePaths {
		outputFile, err = os.Create(path)
		if err == nil {
			savedPath = path
			log.Printf("Creating output file at: %s", path)
			break
		}
	}

	if err != nil {
		log.Printf("Error creating output file: %v", err)
		http.Error(w, "Internal server error - cannot create output file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(converted)
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Created output file: %s", savedPath)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	response := "Результат конвертации:\n\n" + converted + "\n\n"
	response += "Файл сохранен как: " + filepath.Base(savedPath)

	w.Write([]byte(response))
}
