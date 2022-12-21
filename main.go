package main

import (
	"log"
	"net/http"

	"token_generator/handler"
)

func main() {
	authCode := make(chan string)
	th := handler.New(authCode)

	go func() {
		http.HandleFunc("/callback-gl", th.CallbackHandler)

		log.Fatal(http.ListenAndServe(":7000", nil))
	}()

	th.GenerateToken(authCode)
}
