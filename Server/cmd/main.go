package main

import (
	"log"

	"github.com/ArteShow/Assistant/Server/application"
	"github.com/ArteShow/Assistant/Server/internal"
	"github.com/ArteShow/Assistant/Server/pkg/setup"
)

func main() {
	setup.SetUpDatabase()

	go func() {
		log.Println("Starting Internal server...")
		if err := internal.StartInternalServer(); err != nil {
			log.Fatalf("Internal server error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting Internal server...")
		if err := application.StartApplicationServer(); err != nil {
			log.Fatalf("Internal server error: %v", err)
		}
	}()

	select {}
}
