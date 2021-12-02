package main

import (
	"log"
	"net/http"
	//"os/exec"
)

func main() {
	http.HandleFunc("/", showURL)
	http.HandleFunc("/create", create_service)
	http.HandleFunc("/delete", delete_service)
	http.HandleFunc("/download",download_service)
	log.Println("start litening 0.0.0.0:8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", nil))
}
