package main

import (
	"log"
	"os"
	
	"go1fl-sprint6-final-tpl/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "MORSE: ", log.LstdFlags|log.Lshortfile)
	srv := server.NewServer(logger)
	
	logger.Println("Starting Morse Code Converter Service")
	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed to start: ", err)
	}
}