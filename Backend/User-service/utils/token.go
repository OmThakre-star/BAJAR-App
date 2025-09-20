package utils

import (
"errors"
"fmt"
"time"


"github.com/golang-jwt/jwt/v4"
)


var JWTSecret []byte


func InitJWT(secret string) {
JWTSecret = []byte(secret)
}


func CreateToken(userID uint, ttl time.Duration) (string, error) {
tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
Subject: fmt.Sprint(userID),
ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
IssuedAt: jwt.NewNumericDate(time.Now()),
})
return tok.SignedString(JWTSecret)
}


func ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
tok, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
return nil, errors.New("unexpected signing method")
}
return JWTSecret, nil
})
if err != nil {
return nil, err
}
if claims, ok := tok.Claims.(*jwt.RegisteredClaims); ok && tok.Valid {
return claims, nil
}
return nil, errors.New("invalid token")
}