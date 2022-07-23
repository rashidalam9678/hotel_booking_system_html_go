package main

import (
	"fmt"
	"log"
	"net/http"
	
)
const PortNumber= ":8080"



func main(){
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println("Started Serven on: http://localhost:8080 ")
	err := http.ListenAndServe(PortNumber,nil)
	log.Fatal(err)
}