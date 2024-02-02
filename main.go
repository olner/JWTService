package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/createTokens", createTokensHttp)
	http.HandleFunc("/refreshTokens", refreshTokensHttp)

	log.Println("server started on port 80")
	http.ListenAndServe(":80", nil)
}
