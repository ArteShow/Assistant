package main

import (
	"log"

	"github.com/ArteShow/Assistant/Server/application"
)

func main(){
	log.Println("Starting Application server...")
	application.StartServer()
}