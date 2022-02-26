package main 

import (
	"net/http"
	"log"
	"fmt"
	"time"
)

func home(response http.ResponseWriter, request *http.Request){
	fmt.Fprintf(response, "Welcome to home");
	fmt.Println("Hit: Home");
}

func someHeavyRequest(response http.ResponseWriter, request *http.Request){
	time.Sleep(time.Duration(5) * time.Second)
	fmt.Fprintf(response, "Heavy request finished");
	fmt.Println("Heavy request computed")
}

func main(){
	http.HandleFunc("/", home)
	http.HandleFunc("/compute", someHeavyRequest)
	fmt.Println("Started server")
	log.Fatal(http.ListenAndServe(":8000", nil))
}