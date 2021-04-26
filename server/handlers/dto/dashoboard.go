package dto

import "github.com/dgrijalva/jwt-go"

type ConnectionRequest struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type ConnectionAuthorizationResponse struct {
	AppId       string `json:"app_id"`
	ApiKey      string `json:"api_key"`
	AccessToken string `json:"access_token"`
}

type ApplicationResponse struct {
	AppId   string `json:"app_id"`
	AppName string `json:"app_name"`
	AppKey  string `json:"app_key"`
}

type ApplicationAuthorizationClaims struct {
	jwt.StandardClaims
	AppId string
}
