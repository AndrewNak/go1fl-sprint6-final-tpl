package main

import (
	"log"
	"os"
	
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "MORSE: ", log.LstdFlags|log.Lshortfile)
	srv := server.NewServer(logger)
	
	logger.Println("Starting Morse Code Converter Service")
	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed to start: ", err)
	}
}