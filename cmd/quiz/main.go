package main

import (
	"fmt"
	// "time"

	"github.com/jingleMyBells/yapoki-http-api/internal"
	// "database/sql"
	"net/http"
	"log"
)


func main() {
	fmt.Println("Что-то вообще происходит?")

	mux := http.NewServeMux()
	mux.HandleFunc("/", internal.RootHandler)
	mux.HandleFunc("/lalala", internal.PathHandler)
	mux.HandleFunc("/variants/", internal.VariantsHandler)
	mux.HandleFunc("/testing/", internal.TestingHandler)
	mux.HandleFunc("/testing/result/", internal.ResultHandler)


	log.Fatal(http.ListenAndServe(":8080", mux))

	fmt.Println("Похоже закончили")

}
