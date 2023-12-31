package main

import (
	"log"
	"net/http"

	"github.com/jingleMyBells/yapoki-http-api/internal"
)


func main() {
	log.Printf("Запуск")

	err := internal.AddMockData()
	if err != nil {
		log.Printf("Не удалось создать тестовые данные: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", internal.AuthMiddleware(internal.RootHandler))
	mux.HandleFunc("/testing/", internal.TestHTMLHandler)
	mux.HandleFunc("/testing/result/", internal.TestResultHandler)
	mux.HandleFunc("/login/", internal.LoginPageHandler)
	mux.HandleFunc("/api/login/", internal.LoginHandler)
	mux.HandleFunc("/api/variants/", internal.VariantsHandler)
	mux.HandleFunc("/api/testing/", internal.TestingHandler)
	mux.HandleFunc("/api/testing/result/", internal.ResultHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

