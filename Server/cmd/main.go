package main

import (
	"log"

	"github.com/ArteShow/HomeAssistant/Server/application"
)

func main(){
	log.Println("Starting Application server...")
	go func(){
		application.StartServer()
	}()
}