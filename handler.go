package main

import (
	"encoding/json"
	"net/http"
)

func createTokensHttp(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	accessToken, refreshToken := createTokens(guid)
	saveRefreshToken(guid, refreshToken)
	jsonAccessAndRefreshTokens, _ := json.Marshal(map[string]string{"accessToken": accessToken, "refreshToken": refreshToken})
	w.Write(jsonAccessAndRefreshTokens)
}

func refreshTokensHttp(w http.ResponseWriter, r *http.Request) {
	oldRefreshToken := r.URL.Query().Get("refreshToken")
	guid := r.URL.Query().Get("guid")

	if isValidateRefreshToken(guid, oldRefreshToken) {
		newAccessToken, newRefreshToken := createTokens(guid)
		updateRefreshToken(guid, newRefreshToken)
		jsonAccessAndRefreshTokens, _ := json.Marshal(map[string]string{"accessToken": newAccessToken, "refreshToken": newAccessToken})
		w.Write((jsonAccessAndRefreshTokens))
	}
}
