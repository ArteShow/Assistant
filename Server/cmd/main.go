package main

import (
	"log"

	"github.com/ArteShow/Assistant/Server/application"
	"github.com/ArteShow/Assistant/Server/internal"
)

func main(){
	log.Println("Starting Internal server...")
	go func() {
		internal.StartServer()
	}()
	log.Println("Starting Application server...")
	application.StartServer()
}
