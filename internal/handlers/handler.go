package handler

import (
	"encoding/json"
	"net/http"

	db "jwtService/internal/db"
	service "jwtService/internal/services"
)

func CreateTokensHttp(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	accessToken, refreshToken := service.CreateTokens(guid)
	db.SaveRefreshToken(guid, refreshToken)
	jsonAccessAndRefreshTokens, _ := json.Marshal(map[string]string{"accessToken": accessToken, "refreshToken": refreshToken})
	w.Write(jsonAccessAndRefreshTokens)
}

func RefreshTokensHttp(w http.ResponseWriter, r *http.Request) {
	oldRefreshToken := r.URL.Query().Get("refreshToken")
	guid := r.URL.Query().Get("guid")

	if db.IsValidateRefreshToken(guid, oldRefreshToken) {
		newAccessToken, newRefreshToken := service.CreateTokens(guid)
		db.UpdateRefreshToken(guid, newRefreshToken)
		jsonAccessAndRefreshTokens, _ := json.Marshal(map[string]string{"accessToken": newAccessToken, "refreshToken": newAccessToken})
		w.Write((jsonAccessAndRefreshTokens))
	}
}
