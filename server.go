package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

func startServer() {
	log.Println("Starting web server on port", port)
	router := getRouter()
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		panic("Server failed to start " + err.Error())
	}
}
