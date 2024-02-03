package main

import (
	"encoding/json"
	"net/http"
)

func createTokens(guid string) (string, string) {
	accessToken, err := createAccessToken(guid)
	if err != nil {
		panic(err)
	}

	refreshToken, err := createRefreshToken(guid)
	if err != nil {
		panic(err)
	}

	saveRefreshToken(guid, refreshToken)

	return accessToken, refreshToken
}
func createTokensHttp(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	accessToken, err := createAccessToken(guid)
	if err != nil {
		panic(err)
	}
	refreshToken, err := createRefreshToken(guid)
	if err != nil {
		panic(err)
	}
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
