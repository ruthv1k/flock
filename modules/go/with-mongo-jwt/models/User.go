package models

import "github.com/golang-jwt/jwt"

type User struct {
	UserId      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string
}

type UserClaims struct {
	UserId      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	jwt.StandardClaims
}
