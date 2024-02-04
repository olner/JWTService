package main

import (
	handler "jwtService/internal/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/createTokens", handler.CreateTokensHttp)
	http.HandleFunc("/refreshTokens", handler.RefreshTokensHttp)

	log.Println("server started on port 80")
	http.ListenAndServe(":80", nil)
}
