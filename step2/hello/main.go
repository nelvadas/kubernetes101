package main

import (
	"fmt"
	"log"
	"net/http"
        "time"
)

func main() {
	http.HandleFunc("/", homeHandler)
        fmt.Println("Server starting on :8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
        
	currentTime := time.Now().Format("2006.01.02 15:04:05")
        fmt.Fprintf(w,"[%s] INFO Docker and Kubernetes 101   \n",currentTime)
	fmt.Printf("[%s] DEBUG Docker and Kubernetes 101   \n", currentTime)
}
